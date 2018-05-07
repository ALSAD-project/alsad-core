package communicator

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpCommunicator struct {
	url url.URL
}

// NewHTTPCommunicator creates a HTTP communicator for the provided url
func NewHTTPCommunicator(urlString string) (Communicator, error) {
	parsed, err := url.Parse(urlString)
	if err != nil {
		return httpCommunicator{}, err
	}

	return &httpCommunicator{
		url: *parsed,
	}, nil

}

func (c httpCommunicator) FetchData() ([]byte, error) {
	client := http.Client{}
	resp, err := client.Get(c.url.String())

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return respData, nil
}

func (c httpCommunicator) SendData(data []byte) ([]byte, error) {
	client := http.Client{}
	resp, err := client.Post(
		c.url.String(),
		"application/octet-stream",
		bytes.NewReader(data),
	)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return respData, nil
}
