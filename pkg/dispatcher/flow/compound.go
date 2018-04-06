package flow

type compoundFlowImpl struct {
	flows []Flow
}

// CompoundFlow defines the interface of a compound flow
type CompoundFlow interface {
	Flow

	// UnderlyingFlows returns the underlying flows of
	// the compound flow
	UnderlyingFlows() []Flow
}

// NewCompoundFlow creates a compound flow of the input flows
func NewCompoundFlow(flows ...Flow) (CompoundFlow, error) {
	return &compoundFlowImpl{flows: flows}, nil
}

func (f compoundFlowImpl) UnderlyingFlows() []Flow {
	return f.flows
}

func (f compoundFlowImpl) Run(errChan chan<- error) error {
	bufferedErrChan := make(chan error, len(f.flows))
	errs := []error{}
	for _, eachFlow := range f.flows {
		if err := eachFlow.Run(bufferedErrChan); err != nil {
			errs = append(errs, err)
		}
	}

	// FIXME: find a better way to return compound errors
	if len(errs) > 0 {
		return errs[0]
	}

	go func(src <-chan error, dest chan<- error) {
		for eachError := range src {
			dest <- eachError
		}
	}(bufferedErrChan, errChan)

	return nil
}

func (f compoundFlowImpl) Stop() error {
	errs := []error{}
	for _, eachFlow := range f.flows {
		if err := eachFlow.Stop(); err != nil {
			errs = append(errs, err)
		}
	}

	// FIXME: find a better way to return compound errors
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (f compoundFlowImpl) GetState() State {
	flowStates := map[State]int{}
	for _, eachFlow := range f.flows {
		eachFlowState := eachFlow.GetState()
		x := flowStates[eachFlowState]
		flowStates[eachFlowState] = x + 1
	}

	for _, eachState := range []State{
		StateError,
		StateRunning,
		StateInit,
	} {
		if flowStates[eachState] > 0 {
			return eachState
		}
	}

	return StateReady
}
