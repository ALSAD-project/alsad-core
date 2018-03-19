package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/config"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	dpConfig := config.Config{}
	if err := envconfig.Process("dp", &dpConfig); err != nil {
		log.Fatalf("Error on processing configuration: %s", err.Error())
		return
	}

	mux := http.NewServeMux()
	// TODO: add api handlers
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		http.NotFound(resp, req)
	})

	log.Printf("Server is listening on port %d", dpConfig.Port)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", dpConfig.Port),
		mux,
	))
}
