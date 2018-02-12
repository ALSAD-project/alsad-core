package main

import (
	"log"

	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/config"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	dpConfig := config.Config{}
	if err := envconfig.Process("dp", &dpConfig); err != nil {
		log.Fatalf("Error on processing configuration - %s", err.Error())
	}

	log.Printf("Config: %v", dpConfig)
}
