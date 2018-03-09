package expertsystem

// ExpertData for read_data
// ExpertData is the data to be exchanged through HTTP between client and daemon
//
// TODO: Enhance the below problem:
// (Too Lazy to handle different struct, so all data is included and every request)
//
// Pagination is for further rule-based / web-based development
// ---------------------------
type ExpertData struct {
	LinePerPage int
	Page        int

	Filename  string   `json:"filename"`
	Labels    []string `json:"labels"`
	LineCount int      `json:"lineCount"`
	Lines     []string `json:"lines"`
}

// ExpertData ends
// ---------------------------

// For Dev Purpose
// func localConfigure() (config Config, err error) {
// 	config.Port = 4000
// 	config.Host = "127.0.0.1"
// 	config.SrcDir = "tmp/input/"
// 	config.DestDir = "tmp/output/"
// 	return config, nil
// }
