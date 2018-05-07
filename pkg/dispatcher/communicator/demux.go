package communicator

import (
	"errors"
	"sync"
)

type demultiplexingSender struct {
	senders []Sender
}

// NewDemultiplexingSender creates a demultiplexing sender that sends data to
// multiple other senders
func NewDemultiplexingSender(senders ...Sender) (Sender, error) {
	if len(senders) == 0 {
		return nil, errors.New("no senders found")
	}

	return demultiplexingSender{
		senders: senders,
	}, nil
}

func (s demultiplexingSender) SendData(data []byte) ([]byte, error) {
	errChan := make(chan error, len(s.senders))
	wg := sync.WaitGroup{}

	for _, eachSender := range s.senders {
		wg.Add(1)
		go func(sdr Sender, d []byte, ec chan<- error, w *sync.WaitGroup) {
			defer w.Done()
			if _, e := sdr.SendData(d); e != nil {
				ec <- e
			}
		}(eachSender, data, errChan, &wg)
	}

	wg.Wait()
	close(errChan)

	errs := []error{}
	for eachError := range errChan {
		errs = append(errs, eachError)
	}

	// FIXME: find a better way to return compound errors
	if len(errs) > 0 {
		return []byte{}, errs[0]
	}

	return []byte{}, nil
}
