package awair

import (
	"encoding/json"
	"io"
	"time"
)

// awair docs https://support.getawair.com/hc/en-us/articles/360049221014-Awair-Element-Local-API-Feature

type LiveData struct {
	Timestamp      time.Time `json:"timestamp"`
	Score          int       `json:"score"`
	DewPoint       float64   `json:"dew_point"`
	Temp           float64   `json:"temp"`
	Humid          float64   `json:"humid"`
	AbsHumid       float64   `json:"abs_humid"`
	CO2            int       `json:"co2"`
	CO2Est         int       `json:"co2_est"`
	CO2EstBaseline int       `json:"co2_est_baseline"`
	VOC            int       `json:"voc"`
	VOCBaseline    int       `json:"voc_baseline"`
	VOCH2Raw       int       `json:"voc_h2_raw"`
	VOCEthanolRaw  int       `json:"voc_ethanol_raw"`
	PM25           int       `json:"pm25"`
	PM10Est        int       `json:"pm10_est"`
}

func parseLiveData(r io.Reader) (LiveData, error) {
	var ld LiveData
	b, err := io.ReadAll(r)
	if err != nil {
		return ld, err
	}

	err = json.Unmarshal(b, &ld)
	if err != nil {
		return ld, err
	}
	return ld, nil
}
