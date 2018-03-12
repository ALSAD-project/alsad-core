package component

import "github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"

type dataProcessingComponentImpl struct {
	name   string
	sender communicator.Sender
}

// NewBasicDataProcessingComponent creates a basic data processing component
func NewBasicDataProcessingComponent(
	name string,
	sender communicator.Sender,
) DataProcessingComponent {
	return &dataProcessingComponentImpl{
		name:   name,
		sender: sender,
	}
}

func (dpc dataProcessingComponentImpl) GetName() string {
	return dpc.name
}

func (dpc dataProcessingComponentImpl) GetSender() communicator.Sender {
	return dpc.sender
}
