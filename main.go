package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	log.Printf("%s v%s-%s\n", AppName, Version, Commit)

	// Load config and run preflight
	config, err := loadConfig()
	if err != nil {
		log.Println("Failed to load config file:", err)
	}

	if err := preflightCheck(config); err != nil {
		log.Fatalln(err)
	}

	// Initialize client (if any)

	// Initialize collector
	collectorOptions := CollectorOptions{
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(NewCollector(collectorOptions))

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
