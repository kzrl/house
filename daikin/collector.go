package daikin

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type DaikinCollector struct {
	baseURL        string
	httpClient     *http.Client
	endpoints      []string
	powerDesc      *prometheus.Desc
	modeDesc       *prometheus.Desc
	fanRateDesc    *prometheus.Desc
	setTempDesc    *prometheus.Desc
	indoorTempDesc *prometheus.Desc
	upDesc         *prometheus.Desc
}

func NewDaikinCollector(baseURL string) *DaikinCollector {
	endpoints := []string{
		baseURL + "/skyfi/aircon/get_control_info",
		baseURL + "/skyfi/aircon/get_sensor_info",
	}
	
	return &DaikinCollector{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		endpoints: endpoints,
		powerDesc: prometheus.NewDesc(
			"daikin_power",
			"Power (1 = on, 0 = off)",
			nil, nil,
		),
		modeDesc: prometheus.NewDesc(
			"daikin_mode",
			"Mode (0=auto, 1=auto, 2=dry, 3=cool, 4=heat, 6=fan, 7=auto)",
			nil, nil,
		),
		fanRateDesc: prometheus.NewDesc(
			"daikin_fan_rate",
			"Fan Rate",
			nil, nil,
		),
		setTempDesc: prometheus.NewDesc(
			"daikin_set_temp",
			"Set Temperature",
			nil, nil,
		),
		indoorTempDesc: prometheus.NewDesc(
			"daikin_indoor_temp",
			"Indoor Temperature",
			nil, nil,
		),
		upDesc: prometheus.NewDesc(
			"daikin_up",
			"Whether the Daikin unit is reachable (1 = up, 0 = down)",
			[]string{"endpoint"}, nil,
		),
	}
}

func (c *DaikinCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.powerDesc
	ch <- c.modeDesc
	ch <- c.fanRateDesc
	ch <- c.setTempDesc
	ch <- c.indoorTempDesc
	ch <- c.upDesc
}

func (c *DaikinCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := make(map[string]float64)
	allUp := true

	for _, endpoint := range c.endpoints {
		data, err := c.fetchData(endpoint)
		if err != nil {
			log.Printf("daikin: failed to fetch data from %s: %v", endpoint, err)
			ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, endpoint)
			allUp = false
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1, endpoint)
		
		for k, v := range data {
			metrics[k] = v
		}
	}

	if !allUp {
		return
	}

	if v, ok := metrics["pow"]; ok {
		ch <- prometheus.MustNewConstMetric(c.powerDesc, prometheus.GaugeValue, v)
	}
	if v, ok := metrics["mode"]; ok {
		ch <- prometheus.MustNewConstMetric(c.modeDesc, prometheus.GaugeValue, v)
	}
	if v, ok := metrics["stemp"]; ok {
		ch <- prometheus.MustNewConstMetric(c.setTempDesc, prometheus.GaugeValue, v)
	}
	if v, ok := metrics["f_rate"]; ok {
		ch <- prometheus.MustNewConstMetric(c.fanRateDesc, prometheus.GaugeValue, v)
	}
	if v, ok := metrics["htemp"]; ok {
		ch <- prometheus.MustNewConstMetric(c.indoorTempDesc, prometheus.GaugeValue, v)
	}
}

func (c *DaikinCollector) fetchData(endpoint string) (map[string]float64, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp.Body)
}

func (c *DaikinCollector) parseResponse(r io.Reader) (map[string]float64, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	parts := strings.Split(string(b), ",")
	
	for _, part := range parts {
		k, v, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		
		switch k {
		case "pow", "mode", "stemp", "f_rate", "htemp":
			result[k] = val
		}
	}
	
	return result, nil
}