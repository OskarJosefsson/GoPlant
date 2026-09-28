package main

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func SendRequest(url string) (*http.Response, error) {
	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	return response, nil
}
