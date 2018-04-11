/*
 * Example GoLang Daemon.
 *
 * As referencing: https://github.com/golergka/go-tcp-echo/blob/master/go-tcp-echo.go
 */

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/ALSAD-project/alsad-core/pkg/expertsystem/expertdata"
	"github.com/ALSAD-project/alsad-core/pkg/expertsystem/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
)

const (
	configPrefix = "es"
)

var err error
var DaemonConfig config.DaemonConfig

var expertInputCommunicator communicator.Communicator
var expertOutputCommunicator communicator.Communicator

var unlabeledData []string
var labeledData []string

// Request handler for expert client
func readDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Accepted new connection.")

	expertData, err := expertdata.GetExpertDataFromQuery(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write([]byte(nil))
		return
	}

	if len(unlabeledData) == 0 {
		unlabeledData, err = expertdata.FetchData(expertInputCommunicator, DaemonConfig.FsqRequestTimeout)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			w.Write([]byte("\n"))
			w.Write([]byte(nil))
			return
		}

		if len(unlabeledData) == 0 {

			expertData.Lines = []string{}
			expertData.LineCount = 0

			responseJSON, err := json.Marshal(expertData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				w.Write([]byte("\n"))
				w.Write([]byte(nil))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		}
	}

	expertData, err = expertdata.GetExpertDataByLines(unlabeledData, expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write([]byte(nil))
		return
	}

	responseJSON, err := json.Marshal(expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write([]byte(nil))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func updateLabelHandler(w http.ResponseWriter, r *http.Request) {
	var expertData expertdata.ExpertData

	err := json.NewDecoder(r.Body).Decode(&expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write([]byte(nil))
		return
	}

	j := 0
	for i := 0; i < len(expertData.Labels); i++ {
		labeledData = append(labeledData, expertData.Labels[j])
		j++
	}

	if len(unlabeledData) == len(labeledData) {

		log.Println(strconv.Itoa(len(labeledData)) + " records is pushed.")

		if resp, err := expertdata.SendData(labeledData, expertOutputCommunicator); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			w.Write(resp)
			w.Write([]byte("\n"))
			w.Write([]byte(nil))
			return
		} else {
			unlabeledData = []string{}
			labeledData = []string{}
		}
	}
}


func test1(w http.ResponseWriter, r *http.Request) {
	str := []string{"a", "b"}
	data, err := expertdata.SendData(str, expertInputCommunicator)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write([]byte(nil))
		return
	}
	w.WriteHeader(200)
	w.Write(data)
	w.Write([]byte("\n"))
	w.Write([]byte(nil))
	return
}

func main() {
	if err := envconfig.Process(configPrefix, &DaemonConfig); err != nil {
		log.Fatal(err.Error())
	}

	port := flag.Int("port", DaemonConfig.DaemonPort, "Port to accept connections on.")
	flag.Parse()

	s := &http.Server{
		Addr:           "0.0.0.0:" + strconv.Itoa(*port),
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	expertInputCommunicator, err = communicator.NewFSRedisQueueCommunicator(
		DaemonConfig.FsqDir,
		DaemonConfig.FsqExpertInputQueue,
		DaemonConfig.FsqRedisAddr,
	)

	expertOutputCommunicator, err = communicator.NewFSRedisQueueCommunicator(
		DaemonConfig.FsqDir,
		DaemonConfig.FsqExpertOutputQueue,
		DaemonConfig.FsqRedisAddr,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/read_data", readDataHandler)
	http.HandleFunc("/update_label", updateLabelHandler)
	http.HandleFunc("/init", test1)

	log.Println("Listening to connections on port", strconv.Itoa(*port))

	log.Fatal(s.ListenAndServe())
}