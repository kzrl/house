package ecowitt

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type EcowittCollector struct {
	baseURL              string
	httpClient           *http.Client
	tempDesc             *prometheus.Desc
	humidityDesc         *prometheus.Desc
	dewPointDesc         *prometheus.Desc
	solarIrradianceDesc  *prometheus.Desc
	uvIndexDesc          *prometheus.Desc
	tempFeelsLikeDesc    *prometheus.Desc
	windSpeedDesc        *prometheus.Desc
	gustSpeedDesc        *prometheus.Desc
	dayWindMaxDesc       *prometheus.Desc
	soilMoistureDesc     *prometheus.Desc
	batteryDesc          *prometheus.Desc
	upDesc               *prometheus.Desc
}

func NewEcowittCollector(baseURL string) *EcowittCollector {
	return &EcowittCollector{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tempDesc: prometheus.NewDesc(
			"ecowitt_temperature_celsius",
			"Temperature",
			[]string{"location"}, nil,
		),
		humidityDesc: prometheus.NewDesc(
			"ecowitt_humidity",
			"Humidity",
			[]string{"location"}, nil,
		),
		dewPointDesc: prometheus.NewDesc(
			"ecowitt_dew_point_celsius",
			"Dew Point",
			[]string{"location"}, nil,
		),
		solarIrradianceDesc: prometheus.NewDesc(
			"ecowitt_solar_irradiance",
			"Solar irradiance W/m^2",
			nil, nil,
		),
		uvIndexDesc: prometheus.NewDesc(
			"ecowitt_uv_index",
			"UV index",
			nil, nil,
		),
		tempFeelsLikeDesc: prometheus.NewDesc(
			"ecowitt_temperature_feels_like_celsius",
			"Temperature feels like",
			nil, nil,
		),
		windSpeedDesc: prometheus.NewDesc(
			"ecowitt_wind_speed",
			"Wind speed m/s",
			nil, nil,
		),
		gustSpeedDesc: prometheus.NewDesc(
			"ecowitt_gust_speed",
			"Gust speed m/s",
			nil, nil,
		),
		dayWindMaxDesc: prometheus.NewDesc(
			"ecowitt_day_max_wind_speed",
			"Day Max wind speed m/s",
			nil, nil,
		),
		soilMoistureDesc: prometheus.NewDesc(
			"ecowitt_soil_moisture",
			"Soil moisture percentage",
			[]string{"channel", "name"}, nil,
		),
		batteryDesc: prometheus.NewDesc(
			"ecowitt_battery",
			"Battery level",
			[]string{"sensor_type", "channel", "name"}, nil,
		),
		upDesc: prometheus.NewDesc(
			"ecowitt_up",
			"Whether the Ecowitt station is reachable (1 = up, 0 = down)",
			nil, nil,
		),
	}
}

func (c *EcowittCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tempDesc
	ch <- c.humidityDesc
	ch <- c.dewPointDesc
	ch <- c.solarIrradianceDesc
	ch <- c.uvIndexDesc
	ch <- c.tempFeelsLikeDesc
	ch <- c.windSpeedDesc
	ch <- c.gustSpeedDesc
	ch <- c.dayWindMaxDesc
	ch <- c.soilMoistureDesc
	ch <- c.batteryDesc
	ch <- c.upDesc
}

func (c *EcowittCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := c.fetchLiveData()
	if err != nil {
		log.Printf("ecowitt: failed to fetch data: %v", err)
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1)

	c.collectWH25(ch, data.WH25)
	c.collectChAisle(ch, data.ChAisle)
	c.collectChSoil(ch, data.ChSoil)
	c.collectCommonList(ch, data.CommonList)
}

func (c *EcowittCollector) collectWH25(ch chan<- prometheus.Metric, units []WH25) {
	for _, unit := range units {
		if temp, err := strconv.ParseFloat(unit.Temp, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.tempDesc, prometheus.GaugeValue, temp, "office")
		}

		if humidity, err := HumidityToFloat64(unit.Humidity); err == nil {
			ch <- prometheus.MustNewConstMetric(c.humidityDesc, prometheus.GaugeValue, humidity, "office")
		}
	}
}

func (c *EcowittCollector) collectChAisle(ch chan<- prometheus.Metric, channels []ChAisle) {
	for _, channel := range channels {
		if temp, err := TempToFloat64(channel.Temp); err == nil {
			ch <- prometheus.MustNewConstMetric(c.tempDesc, prometheus.GaugeValue, temp, channel.Name)
		}

		if humidity, err := HumidityToFloat64(channel.Humidity); err == nil {
			ch <- prometheus.MustNewConstMetric(c.humidityDesc, prometheus.GaugeValue, humidity, channel.Name)
			
			if temp, err := TempToFloat64(channel.Temp); err == nil {
				if dewpoint, err := DewPoint(temp, humidity); err == nil {
					ch <- prometheus.MustNewConstMetric(c.dewPointDesc, prometheus.GaugeValue, dewpoint, channel.Name)
				}
			}
		}

		if battery, err := strconv.ParseFloat(channel.Battery, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.batteryDesc, prometheus.GaugeValue, battery, "aisle", channel.Channel, channel.Name)
		}
	}
}

func (c *EcowittCollector) collectChSoil(ch chan<- prometheus.Metric, channels []ChSoil) {
	for _, channel := range channels {
		if humidity, err := HumidityToFloat64(channel.Humidity); err == nil {
			ch <- prometheus.MustNewConstMetric(c.soilMoistureDesc, prometheus.GaugeValue, humidity, channel.Channel, channel.Name)
		}

		if battery, err := strconv.ParseFloat(channel.Battery, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.batteryDesc, prometheus.GaugeValue, battery, "soil", channel.Channel, channel.Name)
		}
	}
}

func (c *EcowittCollector) collectCommonList(ch chan<- prometheus.Metric, items []CommonList) {
	var outdoorTemp, outdoorHumidity float64
	
	for _, item := range items {
		switch item.ID {
		case "0x02":
			if temp, err := strconv.ParseFloat(item.Val, 64); err == nil {
				ch <- prometheus.MustNewConstMetric(c.tempDesc, prometheus.GaugeValue, temp, "outdoor")
				outdoorTemp = temp
			}
		case "0x07":
			if humidity, err := HumidityToFloat64(item.Val); err == nil {
				ch <- prometheus.MustNewConstMetric(c.humidityDesc, prometheus.GaugeValue, humidity, "outdoor")
				outdoorHumidity = humidity
			}
		case "3":
			if temp, err := strconv.ParseFloat(item.Val, 64); err == nil {
				ch <- prometheus.MustNewConstMetric(c.tempFeelsLikeDesc, prometheus.GaugeValue, temp)
			}
		case "0x15":
			if irradiance, err := IrradianceToFloat64(item.Val); err == nil {
				ch <- prometheus.MustNewConstMetric(c.solarIrradianceDesc, prometheus.GaugeValue, irradiance)
			}
		case "0x17":
			if uv, err := strconv.ParseFloat(item.Val, 64); err == nil {
				ch <- prometheus.MustNewConstMetric(c.uvIndexDesc, prometheus.GaugeValue, uv)
			}
		case "0x0B":
			if speed, err := SpeedToFloat64(item.Val); err == nil {
				ch <- prometheus.MustNewConstMetric(c.windSpeedDesc, prometheus.GaugeValue, speed)
			}
		case "0x0C":
			if speed, err := SpeedToFloat64(item.Val); err == nil {
				ch <- prometheus.MustNewConstMetric(c.gustSpeedDesc, prometheus.GaugeValue, speed)
			}
		case "0x19":
			if speed, err := SpeedToFloat64(item.Val); err == nil {
				ch <- prometheus.MustNewConstMetric(c.dayWindMaxDesc, prometheus.GaugeValue, speed)
			}
		case "0x03":
			if dewpoint, err := strconv.ParseFloat(item.Val, 64); err == nil {
				ch <- prometheus.MustNewConstMetric(c.dewPointDesc, prometheus.GaugeValue, dewpoint, "outdoor")
			}
		}
	}

	if outdoorHumidity != 0 && outdoorTemp != 0 {
		if dewpoint, err := DewPoint(outdoorTemp, outdoorHumidity); err == nil {
			ch <- prometheus.MustNewConstMetric(c.dewPointDesc, prometheus.GaugeValue, dewpoint, "outdoor-calculated")
		}
	}
}

func (c *EcowittCollector) fetchLiveData() (LiveData, error) {
	url := c.baseURL + "/get_livedata_info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LiveData{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LiveData{}, err
	}
	defer resp.Body.Close()

	return parseLiveData(resp.Body)
}