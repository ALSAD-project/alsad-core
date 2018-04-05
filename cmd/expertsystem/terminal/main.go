/*
 * Example GoLang Daemon.
 *
 * As referencing: https://github.com/golergka/go-tcp-echo/blob/master/go-tcp-echo.go
 */

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ALSAD-project/alsad-core/pkg/expertsystem"
	"github.com/kelseyhightower/envconfig"
)

const (
	configPrefix = "es"
)

type config struct {
	DaemonPort int    `split_words:"true" default:"4000"`
	DaemonHost string `split_words:"true" required:"true"`
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func main() {
	var err error
	var esConfig config
	var expertData expertsystem.ExpertData
	var lineCount int
	var srcFilename string

	if err := envconfig.Process(configPrefix, &esConfig); err != nil {
		log.Fatal(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	port := flag.Int("port", esConfig.DaemonPort, "Port to listen connection on.")
	flag.Parse()

	c := &http.Client{}
	endpoint := "http://" + esConfig.DaemonHost + ":" + strconv.Itoa(*port)

	resp, err := c.Get(endpoint + "/read_data?linePerPage=1&page=0")
	if err != nil {
		log.Fatal(err.Error())
	}

	json.NewDecoder(resp.Body).Decode(&expertData)
	lineCount = expertData.LineCount
	filename, err := newUUID()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print("=============================================\n")
	fmt.Print("ALSAD Expert System Terminal Interface (demo)\n")
	fmt.Print("Please mark label for the below outliers:    \n")
	fmt.Print("---------------------------------------------\n")

	for page := 1; page <= lineCount; page++ {

		// Read data from GET HTTP Request
		resp, err := c.Get(endpoint + "/read_data?linePerPage=1&page=1&srcFilename=" + srcFilename)
		if err != nil {
			log.Fatal(err.Error())
		}
		json.NewDecoder(resp.Body).Decode(&expertData)
		srcFilename = expertData.SrcFilename

		// Ask for user input
		fmt.Print("Outlier [" + strconv.Itoa(page) + "/" + strconv.Itoa(lineCount) + "]: ")
		fmt.Print(expertData.Lines[0] + "\n")
		fmt.Print("Label: ")
		label, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		expertData.Filename = filename
		expertData.Labels = append(expertData.Labels, label)

		// Encode Struct and update daemon via POST HTTP Request
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(expertData)
		resp, err = c.Post(endpoint+"/update_label", "application/json; charset=utf-8", b)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	fmt.Print("All pending outlier is marked.\n")
	fmt.Print("=============================================\n")
}
