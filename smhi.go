package main

import (
	"time"
)

const smhiBaseURL = "https://opendata-download-metobs.smhi.se/api/version/latest"

// Common SMHI parameter IDs. Full list via ListWeatherParameters.
const (
	ParamAirTempHourly     = 1  // °C, instantaneous, hourly
	ParamAirTempDailyMean  = 2  // °C, daily mean
	ParamWindDirection     = 3  // degrees, 10 min mean, hourly
	ParamWindSpeed         = 4  // m/s, 10 min mean, hourly
	ParamPrecipDaily       = 5  // mm, daily sum at 06
	ParamHumidity          = 6  // %, hourly
	ParamPrecipHourly      = 7  // mm, hourly sum
	ParamSnowDepth         = 8  // m, daily at 06
	ParamAirPressure       = 9  // hPa, hourly
	ParamSunshine          = 10 // seconds, hourly sum
	ParamAirTempDailyMin   = 19 // °C
	ParamAirTempDailyMax   = 20 // °C
	ParamGroundState       = 40 // code, daily at 06 (frost etc.)
	ParamPresentWeather    = 13 // code, hourly
	ParamTotalCloudsHourly = 16 // %, hourly
)

// Periods accepted by GetObservations. "corrected-archive" only exists as CSV.
var smhiPeriods = []string{"latest-hour", "latest-day", "latest-months"}

type weatherParameter struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type weatherStation struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Active    bool    `json:"active"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"`
}

type weatherValue struct {
	Date    int64  `json:"date"` // ms epoch, set for instantaneous values
	From    int64  `json:"from"` // ms epoch, set for aggregated values
	To      int64  `json:"to"`
	Ref     string `json:"ref"` // e.g. "2026-10-21" for daily aggregates
	Value   string `json:"value"`
	Quality string `json:"quality"` // G = checked, Y = suspect/unchecked
}

// Time returns the observation time, using the end of the interval for aggregates.
func (v weatherValue) Time() time.Time {
	if v.Date != 0 {
		return time.UnixMilli(v.Date)
	}
	return time.UnixMilli(v.To)
}

type observationParameter struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Unit    string `json:"unit"`
}

type observationStation struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type observationResponse struct {
	Parameter observationParameter `json:"parameter"`
	Station   observationStation   `json:"station"`
	Value     []weatherValue       `json:"value"`
}

type parameterListResponse struct {
	Resource []weatherParameter `json:"resource"`
}

type stationListResponse struct {
	Station []weatherStation `json:"station"`
}

type stationSetEntry struct {
	weatherStation
	Value []weatherValue `json:"value"`
}

type stationSetResponse struct {
	Parameter observationParameter `json:"parameter"`
	Station   []stationSetEntry    `json:"station"`
}
