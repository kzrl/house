package fronius

import (
	"time"
)

type FroniusResponse struct {
	Body Body `json:"Body"`
	Head Head `json:"Head"`
}
type Num1 struct {
	Dt     int     `json:"DT"`
	EDay   float64 `json:"E_Day"`
	ETotal float64 `json:"E_Total"`
	EYear  float64 `json:"E_Year"`
	P      int     `json:"P"`
}
type Inverters struct {
	Num1 Num1 `json:"1"`
}
type Site struct {
	EDay               float64 `json:"E_Day"`
	ETotal             float64 `json:"E_Total"`
	EYear              float64 `json:"E_Year"`
	MeterLocation      string  `json:"Meter_Location"`
	Mode               string  `json:"Mode"`
	PAkku              float64 `json:"P_Akku"`
	PGrid              float64 `json:"P_Grid"`
	PLoad              float64 `json:"P_Load"`
	PPV                float64 `json:"P_PV"`
	RelAutonomy        float64 `json:"rel_Autonomy"`
	RelSelfConsumption float64 `json:"rel_SelfConsumption"`
}
type Data struct {
	Inverters Inverters `json:"Inverters"`
	Site      Site      `json:"Site"`
	Version   string    `json:"Version"`
}
type Body struct {
	Data Data `json:"Data"`
}
type RequestArguments struct {
}
type Status struct {
	Code        int    `json:"Code"`
	Reason      string `json:"Reason"`
	UserMessage string `json:"UserMessage"`
}
type Head struct {
	RequestArguments RequestArguments `json:"RequestArguments"`
	Status           Status           `json:"Status"`
	Timestamp        time.Time        `json:"Timestamp"`
}
