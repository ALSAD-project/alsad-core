package flow

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
)

type basicRateLimitedFlowImpl struct {
	name           string
	dataSource     component.DataSourceComponent
	dataProcessors []component.DataProcessingComponent
	rateLimit      float32

	state      State
	stateMutex sync.RWMutex
}

func newBasicRateLimitedFlow(
	name string,
	dataSource component.DataSourceComponent,
	dataProcessors []component.DataProcessingComponent,
) (RateLimitedFlow, error) {
	return &basicRateLimitedFlowImpl{
		name:           name,
		dataSource:     dataSource,
		dataProcessors: dataProcessors,
		rateLimit:      DefaultRateLimit,
		state:          StateReady,
		stateMutex:     sync.RWMutex{},
	}, nil
}

func (f *basicRateLimitedFlowImpl) GetState() State {
	f.stateMutex.RLock()
	defer f.stateMutex.RUnlock()
	return f.state
}

func (f *basicRateLimitedFlowImpl) SetState(newState State) {
	f.stateMutex.Lock()
	defer f.stateMutex.Unlock()
	f.state = newState
}

func (f *basicRateLimitedFlowImpl) GetRateLimit() float32 {
	return f.rateLimit
}

func (f *basicRateLimitedFlowImpl) SetRateLimit(newRateLimit float32) {
	f.rateLimit = newRateLimit
}

func (f *basicRateLimitedFlowImpl) ensureDispatchingRate(
	last time.Time,
) time.Time {
	minDuration := time.Duration(float32(1*time.Second) / f.rateLimit)
	current := time.Now()
	duration := current.Sub(last)
	if duration < minDuration {
		waitTime := minDuration - duration
		log.Printf(
			"Dispatching rate excesses the limit, continue after %v",
			waitTime,
		)
		<-time.After(waitTime)
		current = current.Add(waitTime)
	}
	return current
}

func (f *basicRateLimitedFlowImpl) dispatchDataSource(
	outputChan chan<- []byte,
	errChan chan<- error,
) {
	lastDispatchTime := time.Time{}
	name := f.dataSource.GetName()
	for {
		// TODO: add stop channel
		lastDispatchTime = f.ensureDispatchingRate(lastDispatchTime)

		log.Printf("Getting data from %s", name)
		data, err := f.dataSource.GetFetcher().FetchData()
		if err != nil {
			log.Printf("Got error from %s: %s", name, err.Error())
			errChan <- err
			continue
		}
		log.Printf("Got %d bytes data from %s.", len(data), name)
		outputChan <- data
	}
}

func (f *basicRateLimitedFlowImpl) dispatchDataDataProcessor(
	dataProcessor component.DataProcessingComponent,
	inputChan <-chan []byte,
	outputChan chan<- []byte,
	errChan chan<- error,
) {
	name := dataProcessor.GetName()
	for inputData := range inputChan {
		log.Printf("Sending %d bytes data to %s.", len(inputData), name)
		outputData, err := dataProcessor.GetSender().SendData(inputData)
		if err != nil {
			log.Printf("Got error from %s: %s", name, err.Error())
			errChan <- err
			continue
		}
		log.Printf("Got %d bytes of data from %s.", len(outputData), name)
		outputChan <- outputData
	}
}

func (f *basicRateLimitedFlowImpl) sinkData(dataSink <-chan []byte) {
	for data := range dataSink {
		log.Printf("Got %d bytes of data from Data sink", len(data))
	}
}

func (f *basicRateLimitedFlowImpl) startDispatch(errChan chan<- error) {
	dataChannels := []chan []byte{}
	for range f.dataProcessors {
		dataChannels = append(dataChannels, make(chan []byte))
	}
	dataSink := make(chan []byte)
	dataChannels = append(dataChannels, dataSink)

	go f.dispatchDataSource(dataChannels[0], errChan)
	for idx, dataProcessor := range f.dataProcessors {
		inputChan := dataChannels[idx]
		outputChan := dataChannels[idx+1]
		go f.dispatchDataDataProcessor(
			dataProcessor,
			inputChan,
			outputChan,
			errChan,
		)
	}

	go f.sinkData(dataSink)
}

func (f *basicRateLimitedFlowImpl) Run(errChan chan<- error) error {
	if f.GetState() != StateReady {
		return fmt.Errorf("cannot start running at %v state", f.state)
	}
	go f.startDispatch(errChan)

	f.SetState(StateRunning)
	return nil
}

func (f *basicRateLimitedFlowImpl) Stop() error {
	return nil
}
