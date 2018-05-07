package main

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/flow"
)

func newFeedbackFlowFromCommunicators(
	expertFetcher communicator.Fetcher,
	slSender communicator.Sender,
	rateLimit float32,
) (flow.Flow, error) {
	f, err := flow.NewBasicRateLimitedFlow(
		"Feedback Flow",
		component.NewBasicDataSourceComponent(
			"Expert Interface Fetcher",
			expertFetcher,
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
		f.SetRateLimit(rateLimit)
	}

	return f, nil
}
