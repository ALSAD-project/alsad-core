package communicator

type memoryCommunicator struct {
	queue chan []byte
}

const (
	// DefaultMemoryQueueSize represents the deault queue size for
	// memory communicator
	DefaultMemoryQueueSize = 20
)

// NewMemoryCommunicator creates a communicator backed by memory
func NewMemoryCommunicator() (Communicator, error) {
	return memoryCommunicator{
		queue: make(chan []byte, DefaultMemoryQueueSize),
	}, nil
}

func (c memoryCommunicator) FetchData() ([]byte, error) {
	d := <-c.queue
	return d, nil
}

func (c memoryCommunicator) SendData(data []byte) ([]byte, error) {
	c.queue <- data
	return []byte{}, nil
}
