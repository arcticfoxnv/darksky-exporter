package main

import (
	"fmt"
	"github.com/arcticfoxnv/darksky-exporter/darksky"
	"github.com/prometheus/client_golang/prometheus"
	forecast "github.com/shawntoffel/darksky"
	"sync"
)

var collectorLabels = []string{
	"latitude",
	"longitude",
	"city",
	"location_name",
}

var (
	apparentTemperatureGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "apparent_temperature",
		Help:      "The apparent (or feels like) temperature in degrees Fahrenheit.",
	}, collectorLabels)

	cloudCoverGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "cloud_cover",
		Help:      "The percentage of sky occluded by clouds, between 0 and 1, inclusive.",
	}, collectorLabels)

	dewPointGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "dew_point",
		Help:      "The dew point in degrees Fahrenheit.",
	}, collectorLabels)

	humidityGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "humidity",
		Help:      "The relative humidity, between 0 and 1, inclusive.",
	}, collectorLabels)

	pressureGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "pressure",
		Help:      "The sea-level air pressure in millibars.",
	}, collectorLabels)

	temperatureGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "temperature",
		Help:      "The air temperature in degrees Fahrenheit.",
	}, collectorLabels)

	windGustGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "wind_gust",
		Help:      "The wind gust speed in miles per hour.",
	}, collectorLabels)

	windSpeedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "darksky",
		Name:      "wind_speed",
		Help:      "The wind speed in miles per hour.",
	}, collectorLabels)
)

type DarkSkyCollectorOptions struct {
	City         string
	Lat          forecast.Measurement
	LocationName string
	Long         forecast.Measurement
}

type DarkSkyCollector struct {
	Options     DarkSkyCollectorOptions
	client      *darksky.Client
	collectLock *sync.Mutex
}

func NewDarkSkyCollector(client *darksky.Client, options DarkSkyCollectorOptions) *DarkSkyCollector {
	return &DarkSkyCollector{
		Options:     options,
		client:      client,
		collectLock: new(sync.Mutex),
	}
}

func (c *DarkSkyCollector) Describe(ch chan<- *prometheus.Desc) {
	apparentTemperatureGauge.Describe(ch)
	cloudCoverGauge.Describe(ch)
	dewPointGauge.Describe(ch)
	humidityGauge.Describe(ch)
	pressureGauge.Describe(ch)
	temperatureGauge.Describe(ch)
	windGustGauge.Describe(ch)
	windSpeedGauge.Describe(ch)
}

func (c *DarkSkyCollector) Collect(ch chan<- prometheus.Metric) {
	c.collectLock.Lock()
	defer c.collectLock.Unlock()

	req := &forecast.ForecastRequest{
		Latitude:  c.Options.Lat,
		Longitude: c.Options.Long,
	}

	data, err := c.client.Forecast(req)

	if err != nil {
		fmt.Printf("Error while getting forecast: %s\n", err)
		return
	}

	labels := make(prometheus.Labels)
	labels["latitude"] = fmt.Sprintf("%f", data.Latitude)
	labels["longitude"] = fmt.Sprintf("%f", data.Longitude)
	labels["city"] = c.Options.City
	labels["location_name"] = c.Options.LocationName

	apparentTemperatureGauge.With(labels).Set(float64(data.Currently.ApparentTemperature))
	cloudCoverGauge.With(labels).Set(float64(data.Currently.CloudCover))
	dewPointGauge.With(labels).Set(float64(data.Currently.DewPoint))
	humidityGauge.With(labels).Set(float64(data.Currently.Humidity))
	pressureGauge.With(labels).Set(float64(data.Currently.Pressure))
	temperatureGauge.With(labels).Set(float64(data.Currently.Temperature))
	windGustGauge.With(labels).Set(float64(data.Currently.WindGust))
	windSpeedGauge.With(labels).Set(float64(data.Currently.WindSpeed))

	apparentTemperatureGauge.Collect(ch)
	cloudCoverGauge.Collect(ch)
	dewPointGauge.Collect(ch)
	humidityGauge.Collect(ch)
	pressureGauge.Collect(ch)
	temperatureGauge.Collect(ch)
	windGustGauge.Collect(ch)
	windSpeedGauge.Collect(ch)
}
