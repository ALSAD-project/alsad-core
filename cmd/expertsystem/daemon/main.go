/*
 * Example GoLang Daemon.
 *
 * As referencing: https://github.com/golergka/go-tcp-echo/blob/master/go-tcp-echo.go
 */

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
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
	"github.com/kelseyhightower/envconfig"
)

const (
	configPrefix = "es"
)

type config struct {
	DaemonPort int    `split_words:"true" default:"4000"`
	SrcDir     string `split_words:"true" required:"true"`
	DestDir    string `split_words:"true" required:"true"`
}

var err error
var esConfig config
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
	srcFilename := query.Get("srcFilename")
	if len(srcFilename) > 0 {
		expertData.SrcFilename = srcFilename
	}

	expertData.Page = page
	expertData.LinePerPage = linePerPage
	return expertData, nil
}

func getLatestFiles(srcDir string) (files []os.FileInfo, err error) {
	files, err = ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() < files[j].ModTime().Unix()
	})

	return files, nil
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

func removeLines(fn string, start int, n int) (err error) {
	if start < 0 {
		return errors.New("invalid request. line numbers start at 1")
	}
	if n < 0 {
		return errors.New("invalid request. negative number to remove")
	}
	var f *os.File
	if f, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return
	}
	defer func() {
		if cErr := f.Close(); err == nil {
			err = cErr
		}
	}()
	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}
	cut, ok := skip(b, start)
	if !ok {
		return errors.New("less than" + strconv.Itoa(start) + "lines")
	}
	if n == -1 {
		return nil
	}
	tail, ok := skip(cut, n)

	if !ok {
		return errors.New("less than" + strconv.Itoa(n) + "lines after line " + strconv.Itoa(start))
	}
	t := int64(len(b) - len(cut))
	if err = f.Truncate(t); err != nil {
		return
	}
	if len(tail) > 0 {
		_, err = f.WriteAt(tail, t)
	} else if t == 0 {
		if err = f.Close(); err != nil {
			return err
		}
		if err = os.Remove(fn); err != nil {
			return err
		}
		return
	}
	return
}

func skip(b []byte, n int) ([]byte, bool) {
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		x := bytes.IndexByte(b, '\n')
		if x < 0 {
			x = len(b)
		} else {
			x++
		}
		b = b[x:]
	}
	return b, true
}

// Request handler for expert client
func readDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Accepted new connection.")
	var files []os.FileInfo

	expertData, err := getExpertData(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	fileName := expertData.SrcFilename
	if len(fileName) == 0 {
		files, err = ioutil.ReadDir(esConfig.SrcDir)
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
		fileName = files[0].Name()
		expertData.SrcFilename = fileName
	}

	srcFile, err := os.Open(esConfig.SrcDir + fileName)
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
	var files []os.FileInfo

	err := json.NewDecoder(r.Body).Decode(&expertData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}

	destFile, err := os.OpenFile(esConfig.DestDir+expertData.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
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
	head := (expertData.Page - 1) * expertData.LinePerPage
	fileName := expertData.SrcFilename
	if len(fileName) == 0 {
		files, err = ioutil.ReadDir(esConfig.SrcDir)
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
		fileName = files[0].Name()
	}

	err = removeLines(esConfig.SrcDir+fileName, head, expertData.LinePerPage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		w.Write([]byte("\n"))
		w.Write(out)
		return
	}
}

func main() {
	if err := envconfig.Process(configPrefix, &esConfig); err != nil {
		log.Fatal(err.Error())
	}

	port := flag.Int("port", esConfig.DaemonPort, "Port to accept connections on.")
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
