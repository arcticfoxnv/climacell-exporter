package main

import (
	"fmt"
	"github.com/arcticfoxnv/climacell-exporter/climacell"
	api "github.com/arcticfoxnv/climacell-go"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var collectorLabels = []string{
	"latitude",
	"longitude",
	"city",
	"location_name",
}

var (
	baroPressureGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "baro_pressure",
		Help:      "Barometric pressure (at surface)",
	}, collectorLabels)
	cloudBaseGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "cloud_base",
		Help:      "The lowest level at which the air contains a perceptible quantity of cloud particles",
	}, collectorLabels)
	cloudCeilingGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "cloud_ceiling",
		Help:      "The height of the lowest layer of clouds which covers more than half of the sky",
	}, collectorLabels)
	cloudCoverGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "cloud_cover",
		Help:      "Fraction of the sky obscured by clouds",
	}, collectorLabels)
	dewpointGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "dewpoint",
		Help:      "Temperature of the dew point",
	}, collectorLabels)
	feelsLikeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "feels_like",
		Help:      "Wind chill and heat window based on season",
	}, collectorLabels)
	humidityGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "humidity",
		Help:      "Percent relative humidity from 0 - 100%",
	}, collectorLabels)
	precipitationGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "precipitation",
		Help:      "Precipitation intensity",
	}, collectorLabels)
	surfaceShortwaveRadiationGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "surface_shortwave_radiation",
		Help:      "Solar radiation reaching the surface",
	}, collectorLabels)
	tempGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Name:      "temp",
		Help:      "Temperature",
	}, collectorLabels)
	visibilityGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "visibility",
		Help:      "Visibility distance",
	}, collectorLabels)
	windDirectionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "wind_direction",
		Help:      "Wind direction in polar degrees 0-360 where 0 is North",
	}, collectorLabels)
	windGustGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "wind_gust",
		Help:      "Wind gust speed",
	}, collectorLabels)
	windSpeedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "climacell",
		Subsystem: "weather",
		Name:      "wind_speed",
		Help:      "Wind speed",
	}, collectorLabels)
)

type CollectorOptions struct {
	City         string
	Latitude     float64
	LocationName string
	Longitude    float64

	EnableWeatherDataLayer bool
}

type Collector struct {
	Options     CollectorOptions
	client      *climacell.Client
	collectLock *sync.Mutex
}

func NewCollector(client *climacell.Client, options CollectorOptions) *Collector {
	return &Collector{
		Options:     options,
		client:      client,
		collectLock: new(sync.Mutex),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	if c.Options.EnableWeatherDataLayer {
		baroPressureGauge.Describe(ch)
		cloudBaseGauge.Describe(ch)
		cloudCeilingGauge.Describe(ch)
		cloudCoverGauge.Describe(ch)
		dewpointGauge.Describe(ch)
		feelsLikeGauge.Describe(ch)
		humidityGauge.Describe(ch)
		precipitationGauge.Describe(ch)
		surfaceShortwaveRadiationGauge.Describe(ch)
		tempGauge.Describe(ch)
		visibilityGauge.Describe(ch)
		windDirectionGauge.Describe(ch)
		windGustGauge.Describe(ch)
		windSpeedGauge.Describe(ch)
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.collectLock.Lock()
	defer c.collectLock.Unlock()

	// Make call to fetch data
	req := &api.RealtimeRequest{
		Latitude:  c.Options.Latitude,
		Longitude: c.Options.Longitude,
		Fields: api.DataFieldList{},
	}

	if c.Options.EnableWeatherDataLayer {
		req.Fields = append(
			req.Fields,
			api.BaroPressure,
			api.CloudBase,
			api.CloudCeiling,
			api.CloudCover,
			api.Dewpoint,
			api.FeelsLike,
			api.Humidity,
			api.Precipitation,
			api.SurfaceShortwaveRadiation,
			api.Temp,
			api.Visibility,
			api.WindDirection,
			api.WindGust,
			api.WindSpeed,
		)
	}

	data, err := c.client.RealtimeWeather(req)
	if err != nil {
		fmt.Printf("Error while getting forecast: %s\n", err)
		return
	}

	// Process data into metrics
	labels := make(prometheus.Labels)
	labels["latitude"] = fmt.Sprintf("%g", data.Latitude)
	labels["longitude"] = fmt.Sprintf("%g", data.Longitude)
	labels["city"] = c.Options.City
	labels["location_name"] = c.Options.LocationName

	c.setIfPresent(baroPressureGauge, labels, data.BaroPressure)
	c.setIfPresent(cloudBaseGauge, labels, data.CloudBase)
	c.setIfPresent(cloudCeilingGauge, labels, data.CloudCeiling)
	c.setIfPresent(cloudCoverGauge, labels, data.CloudCover)
	c.setIfPresent(dewpointGauge, labels, data.Dewpoint)
	c.setIfPresent(feelsLikeGauge, labels, data.FeelsLike)
	c.setIfPresent(humidityGauge, labels, data.Humidity)
	c.setIfPresent(precipitationGauge, labels, data.Precipitation)
	c.setIfPresent(surfaceShortwaveRadiationGauge, labels, data.SurfaceShortwaveRadiation)
	c.setIfPresent(tempGauge, labels, data.Temp)
	c.setIfPresent(visibilityGauge, labels, data.Visibility)
	c.setIfPresent(windDirectionGauge, labels, data.WindDirection)
	c.setIfPresent(windGustGauge, labels, data.WindGust)
	c.setIfPresent(windSpeedGauge, labels, data.WindSpeed)

	// Collect metrics
	c.collectIfPresent(ch, baroPressureGauge, data.BaroPressure)
	c.collectIfPresent(ch, cloudBaseGauge, data.CloudBase)
	c.collectIfPresent(ch, cloudCeilingGauge, data.CloudCeiling)
	c.collectIfPresent(ch, cloudCoverGauge, data.CloudCover)
	c.collectIfPresent(ch, dewpointGauge, data.Dewpoint)
	c.collectIfPresent(ch, feelsLikeGauge, data.FeelsLike)
	c.collectIfPresent(ch, humidityGauge, data.Humidity)
	c.collectIfPresent(ch, precipitationGauge, data.Precipitation)
	c.collectIfPresent(ch, surfaceShortwaveRadiationGauge, data.SurfaceShortwaveRadiation)
	c.collectIfPresent(ch, tempGauge, data.Temp)
	c.collectIfPresent(ch, visibilityGauge, data.Visibility)
	c.collectIfPresent(ch, windDirectionGauge, data.WindDirection)
	c.collectIfPresent(ch, windGustGauge, data.WindGust)
	c.collectIfPresent(ch, windSpeedGauge, data.WindSpeed)
}

func (c *Collector) setIfPresent(gauge *prometheus.GaugeVec, labels prometheus.Labels, data api.DataPoint) {
	if data.Present() {
		gauge.With(labels).Set(data.Value.(float64))
	}
}

func (c *Collector) collectIfPresent(ch chan<- prometheus.Metric, gauge *prometheus.GaugeVec, data api.DataPoint) {
	if data.Present() {
		gauge.Collect(ch)
	}
}
