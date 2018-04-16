package main

type config struct {
	StreamInURL string `envconfig:"STREAM_IN_URL" default:":9999"`
	StreamOutURL string `envconfig:"STREAM_OUT_URL" default:":8888"`
	UserProgram string `envconfig:"USER_PROGRAM" default:"nc"`
	UserProgramArgs string `envconfig:"USER_PROGRAM_ARGS" default:"localhost 8888"`
}