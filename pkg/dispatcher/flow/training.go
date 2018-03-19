package flow

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
)

// NewTrainingFlow creates a training flow
func NewTrainingFlow(
	feederComm communicator.Fetcher,
	baComm communicator.Sender,
	uslComm communicator.Sender,
	expertComm communicator.Communicator,
	slComm communicator.Sender,
) (Flow, error) {
	upper, err := newBasicRateLimitedFlow(
		"Training Upper Flow",
		component.NewBasicDataSourceComponent("Feeder", feederComm),
		[]component.DataProcessingComponent{
			component.NewBasicDataProcessingComponent(
				"Behavior Analyzer",
				baComm,
			),
			component.NewBasicDataProcessingComponent(
				"Unsupervised Learner",
				uslComm,
			),
			component.NewBasicDataProcessingComponent(
				"Expert Interface Sender",
				expertComm,
			),
		},
	)

	if err != nil {
		return nil, err
	}

	lowerFlow, err := newBasicRateLimitedFlow(
		"Training Lower Flow",
		component.NewBasicDataSourceComponent("Expert Interface Fetcher", expertComm),
		[]component.DataProcessingComponent{
			component.NewBasicDataProcessingComponent(
				"Supervised Learner",
				slComm,
			),
		},
	)

	if err != nil {
		return nil, err
	}

	return newCompoundFlow(upper, lowerFlow)
}
