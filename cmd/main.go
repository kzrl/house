package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kzrl/house/awair"
	"github.com/kzrl/house/daikin"
	"github.com/kzrl/house/ecowitt"
	"github.com/kzrl/house/fronius"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Hallo, mein haus!")
	// I know this is not correct German
	// but my 2.5 year old currently speaks a dialect of English that sounds a bit like German.
	// This amuses me and this is my code for fun.

	// Create non-global registry.
	reg := prometheus.NewRegistry()
	
	// Read environment variables and register collectors conditionally
	ecowittURL := os.Getenv("ECOWITT_URL")
	if ecowittURL != "" {
		log.Printf("Registering Ecowitt collector with URL: %s", ecowittURL)
		reg.MustRegister(ecowitt.NewEcowittCollector(ecowittURL))
	} else {
		log.Println("ECOWITT_URL not set, skipping Ecowitt collector")
	}
	
	// Parse Awair hosts from comma-separated list
	// Format: "host=location" or just "host" (uses host as location)
	awairHosts := os.Getenv("AWAIR_HOSTS")
	if awairHosts != "" {
		devices := make(map[string]string)
		hosts := strings.Split(awairHosts, ",")
		for _, host := range hosts {
			host = strings.TrimSpace(host)
			if host == "" {
				continue
			}
			// Use = as separator to avoid conflicts with port numbers
			if strings.Contains(host, "=") {
				parts := strings.SplitN(host, "=", 2)
				devices[parts[0]] = parts[1]
			} else {
				// Use host as both key and location name
				devices[host] = host
			}
		}
		if len(devices) > 0 {
			log.Printf("Registering Awair collector with %d devices", len(devices))
			reg.MustRegister(awair.NewAwairCollector(devices))
		}
	} else {
		log.Println("AWAIR_HOSTS not set, skipping Awair collector")
	}
	
	daikinURL := os.Getenv("DAIKIN_URL")
	if daikinURL != "" {
		log.Printf("Registering Daikin collector with URL: %s", daikinURL)
		reg.MustRegister(daikin.NewDaikinCollector(daikinURL))
	} else {
		log.Println("DAIKIN_URL not set, skipping Daikin collector")
	}
	
	froniusURL := os.Getenv("FRONIUS_URL")
	if froniusURL != "" {
		log.Printf("Registering Fronius collector with URL: %s", froniusURL)
		reg.MustRegister(fronius.NewFroniusCollector(froniusURL))
	} else {
		log.Println("FRONIUS_URL not set, skipping Fronius collector")
	}
	
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting metrics server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
