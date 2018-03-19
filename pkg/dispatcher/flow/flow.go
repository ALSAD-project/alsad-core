package flow

// State defines the type of the flow state
type State string

// Flow defines the interface of a dispatching flow
type Flow interface {
	// Run triggers the dispatching flow to run.
	//
	// This method is expected to run in an asynchronous way. This method
	// returns an error if facing any immediate errors when setting up the flow
	// to run. Other errors during run time should be pass through the channel.
	Run(chan<- error) error

	// Stop triggers the dispatching flow to stop gracefully.
	//
	// This method is expected to run in an asynchronous way. This method
	// returns an error if facing any immediate errors when stopping the flow.
	Stop() error

	// GetState returns the state of the flow.
	GetState() State
}

// RateLimitedFlow defines a dispatching flow which has a rate limit.
//
// The rate limit is defined as the maximum number of data going through
// the flow in a second.
type RateLimitedFlow interface {
	Flow

	// GetRateLimit return the rate limit of the flow.
	GetRateLimit() float32

	// SetRateLimit updates the rate limit of the flow.
	//
	// Please be reminded that updating the rate limit after running the flow
	// may lead to unexpected behaviour.
	SetRateLimit(newRateLimit float32)
}

const (
	// DefaultRateLimit defines the default rate limit for rate-limited flow.
	DefaultRateLimit = float32(1.5)
)

const (
	// StateInit represents the init state of a flow
	StateInit = State("Init")
	// StateReady represents the ready state of a flow
	StateReady = State("Ready")
	// StateRunning represents the running state of a flow
	StateRunning = State("Running")
	// StateError represents the error state of a flow
	StateError = State("Error")
)
