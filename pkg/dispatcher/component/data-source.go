package component

import "github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"

type dataSourceComponentImpl struct {
	name    string
	fetcher communicator.Fetcher
}

// NewBasicDataSourceComponent creates a basic data source component
func NewBasicDataSourceComponent(
	name string,
	fetcher communicator.Fetcher,
) DataSourceComponent {
	return &dataSourceComponentImpl{
		name:    name,
		fetcher: fetcher,
	}
}

func (dsc dataSourceComponentImpl) GetName() string {
	return dsc.name
}

func (dsc dataSourceComponentImpl) GetFetcher() communicator.Fetcher {
	return dsc.fetcher
}
