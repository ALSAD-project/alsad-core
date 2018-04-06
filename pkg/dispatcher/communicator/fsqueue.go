package communicator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/adjust/rmq"
	"github.com/satori/go.uuid"
)

type fsQueueFilenameProducer struct {
	queue rmq.Queue
}

type fsQueueFilenameConsumer struct {
	queue        rmq.Queue
	pollDuration time.Duration

	startedMutex sync.RWMutex
	started      bool

	filenameChan chan string
	ackChan      chan bool
}

type fsRedisQueueCommunicator struct {
	directory        string
	redisConn        rmq.Connection
	queue            rmq.Queue
	filenameProducer fsQueueFilenameProducer
	filenameConsumer fsQueueFilenameConsumer
}

func NewFSRedisQueueCommunicator(directory, queueName, redisAddress string) (Communicator, error) {
	stat, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(directory, 0755); err != nil {
				return nil, fmt.Errorf(
					"failed to create direcotry \"%s\"",
					directory,
				)
			}
		} else {
			if os.IsPermission(err) {
				return nil, fmt.Errorf(
					"permission denied on accessing \"%s\"",
					directory,
				)
			}
			return nil, err
		}
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("path \"%s\" is not a directory", directory)
	}

	conn := rmq.OpenConnection("fs-queue", "tcp", redisAddress, 0)
	rmq.NewCleaner(conn).Clean()

	queue := conn.OpenQueue(queueName)

	filenameProducer, err := newFSRedisQueueProducer(queue)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create filename producer: %s",
			err.Error(),
		)
	}

	consumer, err := newFSRedisQueueConsumer(queue, time.Second)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create filename consumer: %s",
			err.Error(),
		)
	}

	return &fsRedisQueueCommunicator{
		directory:        directory,
		redisConn:        conn,
		queue:            queue,
		filenameProducer: filenameProducer,
		filenameConsumer: consumer,
	}, nil
}

func (c fsRedisQueueCommunicator) FetchData() ([]byte, error) {
	filename := c.filenameConsumer.GetFileName()
	defer c.filenameConsumer.Ack()

	filePath := path.Join(c.directory, filename)

	log.Printf("Reading data from %s", filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read file: %s", err.Error())
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Warning: Failed to delete file %s", filePath)
	}

	return data, nil
}

func (c fsRedisQueueCommunicator) SendData(data []byte) ([]byte, error) {
	filename := uuid.NewV4().String()

	filePath := path.Join(c.directory, filename)
	log.Printf("Writing %d bytes data to %s", len(data), filePath)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return []byte{}, fmt.Errorf("failed to write file: %s", err.Error())
	}

	c.filenameProducer.PutFileName(filename)

	return []byte{}, nil
}

func newFSRedisQueueConsumer(queue rmq.Queue, pollDuration time.Duration) (fsQueueFilenameConsumer, error) {
	return fsQueueFilenameConsumer{
		queue:        queue,
		pollDuration: pollDuration,

		startedMutex: sync.RWMutex{},

		filenameChan: make(chan string),
		ackChan:      make(chan bool),
	}, nil
}

func (c fsQueueFilenameConsumer) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()
	c.filenameChan <- payload
	<-c.ackChan
	delivery.Ack()
}

func (c fsQueueFilenameConsumer) GetFileName() string {
	if !c.IsStartConsuming() {
		c.StartConsuming()
	}

	filename := <-c.filenameChan
	return filename
}

func (c fsQueueFilenameConsumer) IsStartConsuming() bool {
	c.startedMutex.RLock()
	defer c.startedMutex.RUnlock()

	return c.started
}

func (c fsQueueFilenameConsumer) StartConsuming() {
	c.startedMutex.Lock()
	defer c.startedMutex.Unlock()

	if c.started {
		return
	}

	c.started = true
	c.queue.StartConsuming(1, c.pollDuration)
	c.queue.AddConsumer("fs-queue-consumer", c)
}

func (c fsQueueFilenameConsumer) Ack() {
	c.ackChan <- true
}

func newFSRedisQueueProducer(queue rmq.Queue) (fsQueueFilenameProducer, error) {
	return fsQueueFilenameProducer{
		queue: queue,
	}, nil
}

func (p fsQueueFilenameProducer) PutFileName(filename string) {
	p.queue.Publish(filename)
}
