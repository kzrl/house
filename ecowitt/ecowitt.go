package ecowitt

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
)

type LiveData struct {
	CommonList []CommonList `json:"common_list"`
	PiezoRain  []PiezoRain  `json:"piezoRain"`
	WH25       []WH25       `json:"wh25"`
	ChAisle    []ChAisle    `json:"ch_aisle"`
	ChSoil     []ChSoil     `json:"ch_soil"`
}
type WH25 struct {
	Temp     string `json:"intemp"`
	Unit     string `json:"unit"`
	Humidity string `json:"inhumi"`
	Abs      string `json:"abs"`
	Rel      string `json:"rel"`
}

type ChAisle struct {
	Channel  string `json:"channel"`
	Name     string `json:"name"`
	Battery  string `json:"battery"`
	Temp     string `json:"temp"`
	Unit     string `json:"unit"`
	Humidity string `json:"humidity"`
}

// Ch1 big left 3m planter
// Ch2 little round planter
// Ch3 not placed yet
type ChSoil struct { //TODO
	Channel  string `json:"channel"`
	Name     string `json:"name"`
	Battery  string `json:"battery"`
	Humidity string `json:"humidity"`
}

type CommonList struct {
	ID   string `json:"id"`
	Val  string `json:"val"`
	Unit string `json:"unit,omitempty"`
}
type PiezoRain struct {
	ID      string `json:"id"`
	Val     string `json:"val"`
	Battery string `json:"battery,omitempty"`
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


// HumidityToFloat64 strips the % symbol and converts the string to float64.
func HumidityToFloat64(h string) (float64, error) {
	humidity := strings.Replace(h, "%", "", 1)
	return strconv.ParseFloat(humidity, 64)
}

func IrradianceToFloat64(ir string) (float64, error) {
	//stripUnits := strings.Replace(ir, "W/m2", "", 1)
	irradiance := strings.Trim(ir, " W/m2")
	return strconv.ParseFloat(irradiance, 64)
}

func SpeedToFloat64(s string) (float64, error) {
	speed := strings.Trim(s, " m/s")
	return strconv.ParseFloat(speed, 64)
}

func TempToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Function lifted from https://github.com/de-wax/go-pkg/blob/8b504691eda6/dewpoint/dewpoint.go#L16
func DewPoint(T float64, H float64) (float64, error) {
	// Check if the transferred value for the temperature is within the valid range
	if T < -45 || T > 60 {
		return 0, errors.New("Temperatur nicht im gültigen Bereich (-45 - +60°C)")
	} else {

		// Check if the transferred value for humidity is within the valid range
		if H < 0 || H > 100 {
			return 0, errors.New("Feuchtigkeit nicht im gültigen Bereich (0 - 100%)")
		} else {
			// Constants for the Magnus formula
			const a float64 = 17.62
			const b float64 = 243.12
			// Magnus formula
			alpha := math.Log(H/100) + a*T/(b+T)
			return math.Round(((b*alpha)/(a-alpha))*100) / 100, nil
		}
	}
}

