package communicator

type combineCommunicator struct {
	fetcher Fetcher
	sender  Sender
}

func (c combineCommunicator) FetchData() ([]byte, error) {
	return c.fetcher.FetchData()
}

func (c combineCommunicator) SendData(data []byte) ([]byte, error) {
	return c.sender.SendData(data)
}

// NewCombineCommunicator returns a communicator combined from
// a fetcher and a sender
func NewCombineCommunicator(fetcher Fetcher, sender Sender) (Communicator, error) {
	return combineCommunicator{
		fetcher: fetcher,
		sender:  sender,
	}, nil
}
