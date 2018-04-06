package main

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/flow"
)

func newDetectFlowFromCommunicators(
	feederComm communicator.Fetcher,
	baComm communicator.Sender,
	uslComm communicator.Sender,
	slComm communicator.Sender,
	expertComm communicator.Sender,
	rateLimit float32,
) (flow.Flow, error) {
	memoryPipe, err := communicator.NewMemoryCommunicator()
	if err != nil {
		return nil, err
	}

	memoryAndExpertDemux, err := communicator.NewDemultiplexingSender(
		memoryPipe,
		expertComm,
	)
	if err != nil {
		return nil, err
	}

	upper, err := flow.NewBasicRateLimitedFlow(
		"Detect Upper Flow",
		component.NewBasicDataSourceComponent(
			"Feeder",
			feederComm,
		),
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
				"Memory Pipe & Expert Sender Demux",
				memoryAndExpertDemux,
			),
		},
	)
	if err != nil {
		return nil, err
	}

	lower, err := flow.NewBasicRateLimitedFlow(
		"Detect Flow Lower",
		component.NewBasicDataSourceComponent(
			"Memory Pipe",
			memoryPipe,
		),
		[]component.DataProcessingComponent{
			component.NewBasicDataProcessingComponent(
				"Supervised Leaner",
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
