package main

import (
	"fmt"
	"github.com/arcticfoxnv/climacell-exporter/climacell"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
)

func main() {
	log.Printf("%s v%s-%s\n", AppName, Version, Commit)
	log.Println("Powered by ClimaCell - https://www.climacell.co/weather-api")

	// Load config and run preflight
	config, err := loadConfig()
	if err != nil {
		log.Println("Failed to load config file:", err)
	}

	if err := preflightCheck(config); err != nil {
		log.Fatalln(err)
	}

	lat, long, err := LookupCityCoords(config.GetString(CFG_CITY))
	if err != nil {
		log.Fatalln("Failed to lookup city:", err)
	}

	// Initialize client (if any)
	client := climacell.NewClient(
		config.GetString(CFG_API_KEY),
		config.GetDuration(CFG_CACHE_TTL),
	)
	client.SetUserAgent(fmt.Sprintf("climacell-exporter/%s (https://github.com/arcticfoxnv/climacell-exporter)", Version))

	// Initialize collector
	collectorOptions := CollectorOptions{
		City:         config.GetString(CFG_CITY),
		Latitude:          lat,
		LocationName: strings.ToLower(config.GetString(CFG_LOCATION_NAME)),
		Longitude:         long,
		EnableWeatherDataLayer: true,
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(NewCollector(client, collectorOptions))

	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetInt(CFG_LISTEN_PORT)),
		Handler: m,
	}

	// Run
	log.Println("Starting HTTP listener on", s.Addr)
	s.ListenAndServe()
}
