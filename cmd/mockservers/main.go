package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AwairLiveData struct {
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

type ChSoil struct {
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

func generateMockAwairData() AwairLiveData {
	return AwairLiveData{
		Timestamp:      time.Now(),
		Score:          85,
		DewPoint:       12.5,
		Temp:           22.3,
		Humid:          55.2,
		AbsHumid:       10.8,
		CO2:            620,
		CO2Est:         615,
		CO2EstBaseline: 400,
		VOC:            150,
		VOCBaseline:    125,
		VOCH2Raw:       28500,
		VOCEthanolRaw:  19200,
		PM25:           8,
		PM10Est:        12,
	}
}

func generateMockLiveData() LiveData {
	return LiveData{
		CommonList: []CommonList{
			{ID: "0x02", Val: "22.5", Unit: "°C"},
			{ID: "0x07", Val: "65%"},
			{ID: "3", Val: "21.8", Unit: "°C"},
			{ID: "0x15", Val: "850.2 W/m2"},
			{ID: "0x17", Val: "7.5"},
			{ID: "0x0B", Val: "3.2 m/s"},
			{ID: "0x0C", Val: "5.1 m/s"},
			{ID: "0x19", Val: "8.3 m/s"},
			{ID: "0x03", Val: "14.3", Unit: "°C"},
		},
		PiezoRain: []PiezoRain{
			{ID: "rain_rate", Val: "0.0", Battery: "100"},
			{ID: "rain_day", Val: "2.5"},
			{ID: "rain_week", Val: "15.2"},
			{ID: "rain_month", Val: "45.8"},
		},
		WH25: []WH25{
			{
				Temp:     "23.2",
				Unit:     "°C",
				Humidity: "58%",
				Abs:      "1013.25",
				Rel:      "1011.50",
			},
		},
		ChAisle: []ChAisle{
			{
				Channel:  "1",
				Name:     "bedroom",
				Battery:  "90",
				Temp:     "21.5",
				Unit:     "°C",
				Humidity: "62%",
			},
			{
				Channel:  "2",
				Name:     "living_room",
				Battery:  "85",
				Temp:     "22.8",
				Unit:     "°C",
				Humidity: "59%",
			},
			{
				Channel:  "3",
				Name:     "basement",
				Battery:  "95",
				Temp:     "18.2",
				Unit:     "°C",
				Humidity: "71%",
			},
		},
		ChSoil: []ChSoil{
			{
				Channel:  "1",
				Name:     "big_left_3m_planter",
				Battery:  "88",
				Humidity: "45%",
			},
			{
				Channel:  "2",
				Name:     "little_round_planter",
				Battery:  "92",
				Humidity: "38%",
			},
			{
				Channel:  "3",
				Name:     "garden_bed",
				Battery:  "80",
				Humidity: "52%",
			},
		},
	}
}

func ecowittHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	
	mockData := generateMockLiveData()
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mockData); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func awairHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	
	mockData := generateMockAwairData()
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mockData); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
		"service": "ecowitt-mock-server",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := ":8081"
	
	http.HandleFunc("/get_livedata_info", ecowittHandler)
	http.HandleFunc("/air-data/latest", awairHandler)
	http.HandleFunc("/health", healthHandler)
	
	fmt.Printf("Starting mock server for home monitoring devices on port %s\n", port)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  - http://localhost%s/get_livedata_info (Ecowitt weather station data)\n", port)
	fmt.Printf("  - http://localhost%s/air-data/latest (Awair air quality data)\n", port)
	fmt.Printf("  - http://localhost%s/health (Health check)\n", port)
	
	log.Printf("Mock server listening on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}