package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var collectorLabels = []string{
	"name",
}

var (
	sampleMetricGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sample",
		Name:      "metric",
		Help:      "A gauge metric representing something.",
	}, collectorLabels)
)

type CollectorOptions struct {
}

type Collector struct {
	Options     CollectorOptions
	collectLock *sync.Mutex
}

func NewCollector(options CollectorOptions) *Collector {
	return &Collector{
		Options:     options,
		collectLock: new(sync.Mutex),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	sampleMetricGauge.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.collectLock.Lock()
	defer c.collectLock.Unlock()

	// Make call to fetch data

	// Process data into metrics
	labels := make(prometheus.Labels)
	labels["name"] = "foo"

	sampleMetricGauge.With(labels).Set(0.123)

	// Collect metrics
	sampleMetricGauge.Collect(ch)
}
