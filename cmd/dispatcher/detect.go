package main

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/flow"
)

func newDetectFlowFromCommunicators(
	feederFetcher communicator.Fetcher,
	baSender communicator.Sender,
	uslSender communicator.Sender,
	slSender communicator.Sender,
	expertSender communicator.Sender,
	rateLimit float32,
) (flow.Flow, error) {
	memoryPipe, err := communicator.NewMemoryCommunicator()
	if err != nil {
		return nil, err
	}

	memoryAndExpertDemux, err := communicator.NewDemultiplexingSender(
		memoryPipe,
		expertSender,
	)
	if err != nil {
		return nil, err
	}

	upper, err := flow.NewBasicRateLimitedFlow(
		"Detect Upper Flow",
		component.NewBasicDataSourceComponent(
			"Feeder",
			feederFetcher,
		),
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
