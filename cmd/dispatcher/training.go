package main

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/flow"
)

func newTrainingFlowFromCommunicators(
	feederFetcher communicator.Fetcher,
	baSender communicator.Sender,
	uslSender communicator.Sender,
	expertComm communicator.Communicator,
	slSender communicator.Sender,

	rateLimit float32,
) (flow.Flow, error) {
	upper, err := flow.NewBasicRateLimitedFlow(
		"Training Upper Flow",
		component.NewBasicDataSourceComponent("Feeder", feederFetcher),
		[]component.DataProcessingComponent{
			component.NewBasicDataProcessingComponent(
				"Behavior Analyzer",
				baSender,
			),
			component.NewBasicDataProcessingComponent(
				"Unsupervised Learner",
				uslSender,
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
				slSender,
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
