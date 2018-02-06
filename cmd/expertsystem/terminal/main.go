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
	"net/http"
	"os"
	"strconv"

	"github.com/ALSAD-project/alsad-core/pkg/expertsystem"
)

func main() {

	var config expertsystem.Config
	var err error
	var lineCount int
	var expertData expertsystem.ExpertData

	config, err = expertsystem.Configure()
	expertsystem.CheckFatal(err)

	reader := bufio.NewReader(os.Stdin)

	port := flag.Int("port", config.Port, "Port to listen connection on.")
	flag.Parse()

	c := &http.Client{}

	resp, err := c.Get("http://localhost:" + strconv.Itoa(*port) + "/read_data?linePerPage=1&page=0")
	expertsystem.CheckFatal(err)
	json.NewDecoder(resp.Body).Decode(&expertData)
	lineCount = expertData.LineCount
	filename, err := expertsystem.NewUUID()
	expertsystem.CheckFatal(err)

	fmt.Print("=============================================\n")
	fmt.Print("ALSAD Expert System Terminal Interface (demo)\n")
	fmt.Print("Please mark label for the below outliers:    \n")
	fmt.Print("---------------------------------------------\n")

	for page := 1; page <= lineCount; page++ {

		// Read data from GET HTTP Request
		resp, err := c.Get("http://localhost:" + strconv.Itoa(*port) + "/read_data?linePerPage=1&page=" + strconv.Itoa(page))
		expertsystem.CheckFatal(err)
		json.NewDecoder(resp.Body).Decode(&expertData)

		// Ask for user input
		fmt.Print("Outlier [" + strconv.Itoa(page) + "/" + strconv.Itoa(lineCount) + "]: ")
		fmt.Print(expertData.Lines[0] + "\n")
		fmt.Print("Label: ")
		label, err := reader.ReadString('\n')
		expertsystem.CheckFatal(err)
		expertData.Filename = filename
		expertData.Labels = append(expertData.Labels, label)

		// Encode Struct and update daemon via POST HTTP Request
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(expertData)
		resp, err = c.Post("http://localhost:"+strconv.Itoa(*port)+"/update_label", "application/json; charset=utf-8", b)
		expertsystem.CheckFatal(err)
	}
	fmt.Print("All pending outlier is marked.\n")
	fmt.Print("=============================================\n")
}
