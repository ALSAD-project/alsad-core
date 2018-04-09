package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"

	"github.com/kelseyhightower/envconfig"
	"github.com/ALSAD-project/alsad-core/pkg/datafeeder"
)


func main() {

	dfConfig := datafeeder.Config{}
	if err := envconfig.Process("df", &dfConfig); err != nil {
		log.Fatalf("Error on processing configuration: %s", err.Error())
		return
	}

	dist, ok := distmv.NewNormal(
		[]float64{dfConfig.DataMean, dfConfig.NoiseMean},
		mat.NewSymDense(2, []float64{dfConfig.DataVar, 0, 0, dfConfig.NoiseVar}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("Error on creating data source")
	}

	mux := http.NewServeMux()
	// TODO: add api handlers
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		respHandler(resp, req, dist)
	})

	log.Printf("Server is listening on port %d", dfConfig.Port)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", dfConfig.Port),
		mux,
	))
}

func respHandler(w http.ResponseWriter, r *http.Request, dist *distmv.Normal) {
	log.Println("Accepted new connection.")

	v := make([]float64, 2)
	v = dist.Rand(v)

	data := strconv.FormatFloat(v[0] + v[1], 'f', 6, 64)
	noise := strconv.FormatFloat(v[1], 'f',6,64)
	response := data + "," + noise

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}