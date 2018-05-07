package config

// Config defines the configuration of a dispatcher
type Config struct {
	Port int `envconfig:"PORT" default:"8000"`

	EnterMode string `split_words:"true" required:"true"`

	BasicRateLimit float32 `split_words:"true"`

	FsqRedisAddr         string `split_words:"true" required:"true"`
	FsqDir               string `split_words:"true" required:"true"`
	FsqExpertInputQueue  string `split_words:"true" required:"true"`
	FsqExpertOutputQueue string `split_words:"true" required:"true"`

	FeederUrl string `split_words:"true" required:"true"`
	BaUrl     string `split_words:"true" required:"true"`
	UslUrl    string `split_words:"true" required:"true"`
	SlUrl     string `split_words:"true" required:"true"`
}
