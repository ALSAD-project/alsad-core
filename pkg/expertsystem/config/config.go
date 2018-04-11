package config


// DaemonConfig defines configuration of an expert system daemon
type DaemonConfig struct {
	DaemonPort int    `split_words:"true" default:"4000"`
	SrcDir     string `split_words:"true" required:"true"`
	DestDir    string `split_words:"true" required:"true"`

	FsqRedisAddr         string `split_words:"true" required:"true"`
	FsqDir               string `split_words:"true" required:"true"`
	FsqExpertInputQueue  string `split_words:"true" required:"true"`
	FsqExpertOutputQueue string `split_words:"true" required:"true"`
	FsqRequestTimeout	int		`split_words:"true" default:"10"`
}

// TerminalConfig defines configuration of an expert system terminal client
type TerminalConfig struct {
	DaemonPort int    `split_words:"true" default:"4000"`
	DaemonHost string `split_words:"true" required:"true"`
}