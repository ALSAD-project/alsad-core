package main

type config struct {
	DispatcherListenURL string `envconfig:"DISPATCHER_LISTEN_URL" default:":9999"`
	UserProgListenURL string `envconfig:"USERPROG_LISTEN_URL" default:":8888"`
	UserProgram string `envconfig:"USER_PROGRAM" default:"nc localhost 8888"`
}