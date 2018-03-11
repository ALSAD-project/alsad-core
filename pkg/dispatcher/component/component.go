package component

import "github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"

// Component defines the common interface of a component for dispatcher
type Component interface {
	// GetName returns the name of the component. The name is usually to
	// identify the component from another.
	GetName() string
}

// DataSourceComponent defines the interface of a component which only
// acts as a data source.
type DataSourceComponent interface {
	Component

	GetFetcher() communicator.Fetcher
}

// DataProcessingComponent defines the interface of a component which
// acts as a data processor.
type DataProcessingComponent interface {
	Component

	GetSender() communicator.Sender
}
