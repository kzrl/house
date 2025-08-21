# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based home monitoring system that collects metrics from various smart home devices and exposes them as Prometheus metrics on port 8080. The application integrates with:
- Awair air quality monitors
- Ecowitt weather stations
- Daikin air conditioners
- Fronius solar inverters (currently commented out)

## Build and Development Commands

```bash
# Build the application
go build -o house cmd/main.go

# Build for Linux deployment
GOOS=linux GOARCH=amd64 go build -o house cmd/main.go

# Run tests
go test ./...

# Run a specific test
go test ./dashboard

# Deploy to production server (192.168.1.8)
./deploy.sh
```

## Architecture

The application follows a modular package structure where each device integration is a separate package:

- **cmd/main.go**: Entry point that initializes all metrics collectors and starts the HTTP server
- **awair/**: Fetches air quality data from Awair Element devices via local API
- **ecowitt/**: Collects weather station data including temperature, humidity, soil moisture
- **daikin/**: Monitors Daikin AC unit status including power, mode, temperatures
- **fronius/**: Solar inverter integration for power generation metrics
- **dashboard/**: Contains utility functions for electricity rate calculations based on time-of-use pricing

Each device package implements the Prometheus Collector pattern:
1. Defines its own data structures for API responses
2. Implements a `collector.go` file with a custom Collector type
3. Provides `Describe()` and `Collect()` methods for on-demand metric collection
4. Fetches metrics only when Prometheus scrapes the endpoint (no background polling)

The application uses a non-global Prometheus registry and exposes all metrics on `/metrics` endpoint for scraping by Prometheus.

## Environment Variables

See `.env.example` for required configuration. All device integrations are optional and only registered if their respective URL environment variables are set.