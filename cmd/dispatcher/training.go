package main

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/flow"
)

func newTrainingFlowFromCommunicators(
	feederComm communicator.Fetcher,
	baComm communicator.Sender,
	uslComm communicator.Sender,
	expertComm communicator.Communicator,
	slComm communicator.Sender,

	rateLimit float32,
) (flow.Flow, error) {
	upper, err := flow.NewBasicRateLimitedFlow(
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

	lower, err := flow.NewBasicRateLimitedFlow(
		"Training Lower Flow",
		component.NewBasicDataSourceComponent(
			"Expert Interface Fetcher",
			expertComm,
		),
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

	if rateLimit > 0 {
		upper.SetRateLimit(rateLimit)
		lower.SetRateLimit(rateLimit)
	}

	return flow.NewCompoundFlow(upper, lower)
}
