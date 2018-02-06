/*
 * Example GoLang Daemon.
 *
 * As referencing: https://github.com/golergka/go-tcp-echo/blob/master/go-tcp-echo.go
 */

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"

	"github.com/ALSAD-project/alsad-core/pkg/expertsystem"
)

var config expertsystem.Config
var err error
var out []byte

func getExpertData(query url.Values) (expertData expertsystem.ExpertData, err error) {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		return expertData, err
	}
	linePerPage, err := strconv.Atoi(query.Get("linePerPage"))
	if err != nil {
		return expertData, err
	}
	expertData.Page = page
	expertData.LinePerPage = linePerPage
	return expertData, nil
}

func readLine(r io.Reader, expertData expertsystem.ExpertData) (resExpertData expertsystem.ExpertData, err error) {
	sc := bufio.NewScanner(r)
	lastLine := 0

	head := (expertData.Page - 1) * expertData.LinePerPage
	tail := expertData.Page * expertData.LinePerPage

	for sc.Scan() {
		if lastLine >= head && lastLine < tail {
			expertData.Lines = append(expertData.Lines, sc.Text())
		}
		lastLine++
	}
	expertData.LineCount = lastLine
	return expertData, nil
}

// Request handler for expert client
func readDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Accepted new connection.")

	expertData, err := getExpertData(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	files, err := ioutil.ReadDir(config.SrcDir)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	if len(files) == 0 {
		w.WriteHeader(200)
		w.Write([]byte("No pending data is to be labeled."))
		w.Write(out)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() < files[j].ModTime().Unix()
	})
	fileName := files[0].Name()

	srcFile, err := os.Open(config.SrcDir + fileName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}
	defer srcFile.Close()

	expertData, err = readLine(bufio.NewReader(srcFile), expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	responseJSON, err := json.Marshal(expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func updateLabelHandler(w http.ResponseWriter, r *http.Request) {
	var expertData expertsystem.ExpertData
	err := json.NewDecoder(r.Body).Decode(&expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	destFile, err := os.OpenFile(config.DestDir+expertData.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}
	writer := bufio.NewWriter(destFile)
	defer destFile.Close()

	for i := 0; i < len(expertData.Lines); i++ {
		updatedData := expertData.Lines[i] + "," + expertData.Labels[i]
		_, err := writer.WriteString(updatedData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			w.Write([]byte("\n"))
			w.Write(out)
			return
		}
		writer.Flush()
	}
}

func main() {
	config, err = expertsystem.Configure()
	expertsystem.CheckFatal(err)

	port := flag.Int("port", config.Port, "Port to accept connections on.")
	flag.Parse()

	s := &http.Server{
		Addr:           "0.0.0.0:" + strconv.Itoa(*port),
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	http.HandleFunc("/read_data", readDataHandler)
	http.HandleFunc("/update_label", updateLabelHandler)

	log.Println("Listening to connections on port", strconv.Itoa(*port))

	log.Fatal(s.ListenAndServe())
}
