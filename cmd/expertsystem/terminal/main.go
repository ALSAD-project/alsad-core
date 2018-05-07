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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ALSAD-project/alsad-core/pkg/expertsystem/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/ALSAD-project/alsad-core/pkg/expertsystem/expertdata"
	"strings"
)

const (
	configPrefix = "es"
)

func main() {
	var esConfig config.TerminalConfig
	var expertData expertdata.ExpertData
	var lineCount int

	if err := envconfig.Process(configPrefix, &esConfig); err != nil {
		log.Fatal(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	port := flag.Int("port", esConfig.DaemonPort, "Port to listen connection on.")
	flag.Parse()

	c := &http.Client{}
	endpoint := "http://" + esConfig.DaemonHost + ":" + strconv.Itoa(*port)

	fmt.Print("=============================================\n")
	fmt.Print("ALSAD Expert System Terminal Interface (demo)\n")
	fmt.Print("Please mark label for the below outliers:    \n")
	fmt.Print("---------------------------------------------\n")


	for {
		resp, err := c.Get(endpoint + "/read_data?linePerPage=0&page=0")
		if err != nil {
			log.Fatal(err.Error())
		}

		json.NewDecoder(resp.Body).Decode(&expertData)
		lineCount = expertData.LineCount

		if lineCount == 0 {
			break
		}

		for page := 1; page <= lineCount; page++ {

			// Read data from GET HTTP Request
			resp, err := c.Get(endpoint + "/read_data?linePerPage=1&page=" + strconv.Itoa(page))
			if err != nil {
				log.Fatal(err.Error())
			}
			json.NewDecoder(resp.Body).Decode(&expertData)

			// Ask for user input
			fmt.Print("Outlier [" + strconv.Itoa(page) + "/" + strconv.Itoa(lineCount) + "]: ")
			fmt.Print(expertData.Lines[0] + "\n")
			fmt.Print("Label: ")
			label, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err.Error())
			}
			label = strings.TrimSuffix(label, "\n")
			expertData.Labels = append(expertData.Labels, expertData.Lines[0]+","+label)
			expertData.Page = page
			expertData.LinePerPage = 1

			// Encode Struct and update daemon via POST HTTP Request
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(expertData)
			resp, err = c.Post(endpoint+"/update_label", "application/json; charset=utf-8", b)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		fmt.Print("Label is updated, searching for next file....\n")
	}
	fmt.Print("All pending outlier is marked.\n")
	fmt.Print("=============================================\n")
}