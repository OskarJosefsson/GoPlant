package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestListStationsReturnsActiveStationsWithBothParameters(t *testing.T) {
	previousClient := httpClient
	httpClient = &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			var body string
			switch request.URL.Path {
			case "/api/version/latest/parameter/1.json":
				body = `{"station":[
					{"id":2,"name":"Beta","active":true},
					{"id":1,"name":"Alpha","active":true},
					{"id":3,"name":"Inactive","active":false}
				]}`
			case "/api/version/latest/parameter/5.json":
				body = `{"station":[
					{"id":1,"name":"Alpha","active":true},
					{"id":2,"name":"Beta","active":false},
					{"id":3,"name":"Inactive","active":true}
				]}`
			default:
				t.Fatalf("unexpected request path: %s", request.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
				Request:    request,
			}, nil
		}),
	}
	t.Cleanup(func() {
		httpClient = previousClient
	})

	stations, err := listStations()
	if err != nil {
		t.Fatalf("listStations() error = %v", err)
	}
	if len(stations.Station) != 1 {
		t.Fatalf("got %d stations, want 1", len(stations.Station))
	}
	if got := stations.Station[0].ID; got != 1 {
		t.Errorf("station ID = %d, want 1", got)
	}
	if got := stations.Station[0].Name; got != "Alpha" {
		t.Errorf("station name = %q, want Alpha", got)
	}
}

func TestLatestStationWeatherFetchesTemperatureAndDailyRainfall(t *testing.T) {
	previousClient := httpClient
	httpClient = &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			var body string
			switch request.URL.Path {
			case "/api/version/latest/parameter/1/station/75250/period/latest-day/data.json":
				body = `{"value":[
					{"date":2000,"value":"11.0"},
					{"date":3000,"value":"12.5"}
				]}`
			case "/api/version/latest/parameter/5/station/75250/period/latest-day/data.json":
				body = `{"value":[{"from":1000,"to":4000,"ref":"2026-10-29","value":"1.2"}]}`
			default:
				t.Fatalf("unexpected request path: %s", request.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
				Request:    request,
			}, nil
		}),
	}
	t.Cleanup(func() {
		httpClient = previousClient
	})

	weather, err := latestStationWeather(75250)
	if err != nil {
		t.Fatalf("latestStationWeather() error = %v", err)
	}
	if got := weather.Temperature.Value; got != "12.5" {
		t.Errorf("temperature = %q, want 12.5", got)
	}
	if got := weather.Rainfall.Value; got != "1.2" {
		t.Errorf("rainfall = %q, want 1.2", got)
	}
	if got := weather.Rainfall.Ref; got != "2026-10-29" {
		t.Errorf("rainfall date = %q, want 2026-10-29", got)
	}
}
