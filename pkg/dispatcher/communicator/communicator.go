package communicator

// Fetcher defines the interface of a fetcher for dispatcher to get data
// from another component.
type Fetcher interface {
	// FetchData fetches data from a component
	FetchData() ([]byte, error)
}

// Sender defines the interface of a sender for dispatcher to send data
// to another component.
type Sender interface {
	// SendData sends data to a component.
	SendData([]byte) ([]byte, error)
}

// Communicator defines the interface of a communicator for dispatcher to
// communicate with another component. Communicator is basically a union
// of Fetcher and Sender
type Communicator interface {
	Fetcher
	Sender
}
