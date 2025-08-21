package awair

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type AwairCollector struct {
	devices     map[string]string
	httpClient  *http.Client
	scoreDesc   *prometheus.Desc
	tempDesc    *prometheus.Desc
	humidDesc   *prometheus.Desc
	co2Desc     *prometheus.Desc
	vocDesc     *prometheus.Desc
	pm25Desc    *prometheus.Desc
	pm10Desc    *prometheus.Desc
	upDesc      *prometheus.Desc
}

func NewAwairCollector(devices map[string]string) *AwairCollector {
	return &AwairCollector{
		devices: devices,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		scoreDesc: prometheus.NewDesc(
			"awair_score",
			"Awair Score",
			[]string{"location"}, nil,
		),
		tempDesc: prometheus.NewDesc(
			"awair_temperature_celsius",
			"Temperature",
			[]string{"location"}, nil,
		),
		humidDesc: prometheus.NewDesc(
			"awair_humidity",
			"Relative Humidity %",
			[]string{"location"}, nil,
		),
		co2Desc: prometheus.NewDesc(
			"awair_co2",
			"CO2 ppm",
			[]string{"location"}, nil,
		),
		vocDesc: prometheus.NewDesc(
			"awair_voc",
			"Volatile Organic Compounds ppb",
			[]string{"location"}, nil,
		),
		pm25Desc: prometheus.NewDesc(
			"awair_pm25",
			"pm2.5 µg/m³",
			[]string{"location"}, nil,
		),
		pm10Desc: prometheus.NewDesc(
			"awair_pm10",
			"pm10 µg/m³",
			[]string{"location"}, nil,
		),
		upDesc: prometheus.NewDesc(
			"awair_up",
			"Whether the Awair device is reachable (1 = up, 0 = down)",
			[]string{"location", "ip"}, nil,
		),
	}
}

func (c *AwairCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.scoreDesc
	ch <- c.tempDesc
	ch <- c.humidDesc
	ch <- c.co2Desc
	ch <- c.vocDesc
	ch <- c.pm25Desc
	ch <- c.pm10Desc
	ch <- c.upDesc
}

func (c *AwairCollector) Collect(ch chan<- prometheus.Metric) {
	for ip, location := range c.devices {
		c.collectDevice(ch, ip, location)
	}
}

func (c *AwairCollector) collectDevice(ch chan<- prometheus.Metric, ip, location string) {
	data, err := c.fetchLiveData(ip)
	if err != nil {
		log.Printf("awair: failed to fetch data from %s (%s): %v", ip, location, err)
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, location, ip)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1, location, ip)
	ch <- prometheus.MustNewConstMetric(c.scoreDesc, prometheus.GaugeValue, float64(data.Score), location)
	ch <- prometheus.MustNewConstMetric(c.tempDesc, prometheus.GaugeValue, data.Temp, location)
	ch <- prometheus.MustNewConstMetric(c.humidDesc, prometheus.GaugeValue, data.Humid, location)
	ch <- prometheus.MustNewConstMetric(c.co2Desc, prometheus.GaugeValue, float64(data.CO2), location)
	ch <- prometheus.MustNewConstMetric(c.vocDesc, prometheus.GaugeValue, float64(data.VOC), location)
	ch <- prometheus.MustNewConstMetric(c.pm25Desc, prometheus.GaugeValue, float64(data.PM25), location)
	ch <- prometheus.MustNewConstMetric(c.pm10Desc, prometheus.GaugeValue, float64(data.PM10Est), location)
}

func (c *AwairCollector) fetchLiveData(hostname string) (LiveData, error) {
	endpoint := fmt.Sprintf("http://%s/air-data/latest", hostname)
	req, err := http.NewRequest("GET", endpoint, nil)
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