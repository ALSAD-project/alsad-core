package config

// Config defines the configuration of a dispatcher
type Config struct {
	Port              int    `envconfig:"PORT" default:"8000"`
	EnterMode         string `split_words:"true" required:"true"`
	ExpertDb          string `split_words:"true" required:"true"`
	ExpertDbInputDir  string `split_words:"true" required:"true"`
	ExpertDbOutputDir string `split_words:"true" required:"true"`
	FeederURL         string `envconfig:"FEEDER_URL" required:"true"`
	BaURL             string `envconfig:"BA_URL" required:"true"`
	UslURL            string `envconfig:"USL_URL" required:"true"`
	SlURL             string `envconfig:"SL_URL" required:"true"`
}
