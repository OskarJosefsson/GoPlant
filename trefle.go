package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const trefleBaseURL = "https://trefle.io/api/v1/plants"

type searchPlant struct {
	ID   int    `json:"id"`
	Name string `json:"common_name"`
	Slug string `json:"slug"`
}

type plantTaxon struct {
	Name string `json:"name"`
}

type fullPlant struct {
	ID             int        `json:"id"`
	CommonName     string     `json:"common_name"`
	ScientificName string     `json:"scientific_name"`
	Family         plantTaxon `json:"family"`
	Genus          plantTaxon `json:"genus"`
	ImageURL       string     `json:"image_url"`
	Status         string     `json:"status"`
	Year           int        `json:"year"`
}

type searchResponse struct {
	Data []searchPlant `json:"data"`
}

type fullPlantResponse struct {
	Data fullPlant `json:"data"`
}

func GetToken() string {
	return os.Getenv("PLANT_API_KEY")
}

func SearchPlants(searchTerm string) (searchResponse, error) {
	token := GetToken()
	if token == "" {
		return searchResponse{}, fmt.Errorf("PLANT_API_KEY is not set")
	}

	query := url.Values{}
	query.Set("token", token)
	query.Set("q", searchTerm)
	requestURL := trefleBaseURL + "/search?" + query.Encode()

	response, err := SendRequest(requestURL)
	if err != nil {
		return searchResponse{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return searchResponse{}, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return searchResponse{}, fmt.Errorf("plant search failed with status %s: %s", response.Status, string(data))
	}

	var search searchResponse
	if err := json.Unmarshal(data, &search); err != nil {
		return searchResponse{}, err
	}

	return search, nil
}

func GetFullPlant(identifier string) (fullPlant, error) {
	token := GetToken()
	if token == "" {
		return fullPlant{}, fmt.Errorf("PLANT_API_KEY is not set")
	}

	query := url.Values{}
	query.Set("token", token)
	requestURL := trefleBaseURL + "/" + identifier + "?" + query.Encode()

	response, err := SendRequest(requestURL)
	if err != nil {
		return fullPlant{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fullPlant{}, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fullPlant{}, fmt.Errorf("plant request failed with status %s: %s", response.Status, string(data))
	}

	var plant fullPlantResponse
	if err := json.Unmarshal(data, &plant); err != nil {
		return fullPlant{}, err
	}

	return plant.Data, nil
}
