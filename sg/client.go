package sg

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	Name string
	Data CompletionEvent
}

type HttpClient struct {
	Client  *http.Client
	Headers map[string]string
}

type CompletionEvent struct {
	Completion string `json:"completion"`
	StopReason string `json:"stopReason"`
}

var httpClient HttpClient

func initClient() error {
	httpClient = HttpClient{
		Client:  &http.Client{},
		Headers: make(map[string]string),
	}

	authToken := os.Getenv("CODY_AUTH_TOKEN")
	if authToken != "" {
		httpClient.AddHeader("Authorization", fmt.Sprintf("token %s", authToken))
	} else {
		return fmt.Errorf("CODY_AUTH_TOKEN not set")
	}

	return nil
}

func (cli *HttpClient) AddHeader(key string, value string) {
	cli.Headers[key] = value
}

func (cli *HttpClient) PostRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range cli.Headers {
		req.Header.Add(key, value)
	}

	resp, err := cli.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
