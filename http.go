package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

const baseUrl = "https://trefle.io/api/v1/plants"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type searchPlant struct {
	ID   int    `json:"id"`
	Name string `json:"common_name"`
}

type searchResponse struct {
	Data []searchPlant `json:"data"`
}

func SendRequest(url string) (*http.Response, error) {
	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetToken() string {
	return os.Getenv("PLANT_API_KEY")
}

func SearchPlants(url string, searchTerm string) (searchResponse, error) {

	token := GetToken()

	if url == "" {
		url = baseUrl + "/search?token=" + token + "&q=" + searchTerm
	}

	response, err := SendRequest(url)
	if err != nil {
		return searchResponse{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return searchResponse{}, err
	}

	var search searchResponse

	if err := json.Unmarshal(data, &search); err != nil {
		return searchResponse{}, err
	}

	return search, nil

}
