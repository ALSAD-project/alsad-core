package flow

import (
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/component"
)

// NewTrainingFlow creates a training flow
func NewTrainingFlow(
	dataSource component.DataSourceComponent,
	behaviorAnalyzer component.DataProcessingComponent,
	unsupervisedLeaner component.DataProcessingComponent,
	expertInterface component.DataProcessingComponent,
	supervisedLeaner component.DataProcessingComponent,
) RateLimitedFlow {
	f := newBasicRateLimitedFlow(
		"training",
		dataSource,
		[]component.DataProcessingComponent{
			behaviorAnalyzer,
			unsupervisedLeaner,
			expertInterface,
			supervisedLeaner,
		},
		DefaultRateLimit,
	)

	f.SetState(StateReady)

	return &f
}
