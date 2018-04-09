package datafeeder

// Config defines the configuration of a dispatcher
type Config struct {
	Port              int    `envconfig:"PORT" default:"8080"`
	DataMean		  float64 `envconfig:"DATA_MEAN" default:"100"`
	DataVar			  float64 `envconfig:"DATA_VAR" default:"20"`
	NoiseMean		  float64 `envconfig:"NOISE_MEAN" default:"0"`
	NoiseVar		  float64 `envconfig:"NOISE_VAR" default:"2"`
}
