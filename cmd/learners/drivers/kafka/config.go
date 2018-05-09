package main

type config struct {
	DispatcherListenURL string `envconfig:"DISPATCHER_LISTEN_URL" default:":9999"`
	KafkaBrokerURL string `envconfig:"KAFKA_BROKER_URL" default:":9093"`
}