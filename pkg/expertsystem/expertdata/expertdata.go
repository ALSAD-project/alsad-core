package expertdata

import (
	"net/url"
	"strconv"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"time"
	"strings"
)

// ExpertData is the data to be exchanged through HTTP between client and daemon
type ExpertData struct {
	LinePerPage int
	Page        int

	Filename    string   `json:"filename"`
	Labels      []string `json:"labels"`
	LineCount   int      `json:"lineCount"`
	Lines       []string `json:"lines"`
}

func GetExpertDataFromQuery(query url.Values) (expertData ExpertData, err error) {
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

func GetExpertDataByLines(unlabeledData []string, expertData ExpertData) (resExpertData ExpertData, err error) {
	expertData.LineCount = len(unlabeledData)

	if expertData.LinePerPage != 0 {
		head := (expertData.Page - 1) * expertData.LinePerPage
		tail := expertData.Page * expertData.LinePerPage

		for i := head; i < tail; i++ {
			expertData.Lines = append(expertData.Lines, unlabeledData[i])
		}
	}
	return expertData, nil
}

func FetchData(expertInputCommunicator communicator.Communicator, requestTimeLimit int) (data[]string, err error){
	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		tempData, tempErr := expertInputCommunicator.FetchData()
		dataChan <- tempData
		errChan <- tempErr
	} ()

	select {

	case srcData := <-dataChan:
		err = <- errChan
		if err != nil {
			return nil, err
		}
		return strings.Split(string(srcData), "\n"), nil

	case <-time.After(time.Duration(requestTimeLimit) * time.Second):
		// call timed out
		return []string{}, nil
	}
}

func SendData(sourceData []string, expertInputCommunicator communicator.Communicator) ([]byte, error){

	data, err := expertInputCommunicator.SendData([]byte(strings.Join(sourceData, "\n")))

	if err != nil {
		return nil, err
	}

	return data, nil
}