package main

import (
	"fmt"
	"github.com/arcticfoxnv/darksky-exporter/darksky"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	forecast "github.com/shawntoffel/darksky"
	"log"
	"net/http"
	"strings"
)

func main() {
	log.Printf("darksky-exporter v%s-%s\n", Version, Commit)
	log.Println("Powered by Dark Sky - https://darksky.net/poweredby/")

	config, err := loadConfig()
	if err != nil {
		log.Println("Failed to load config file:", err)
	}

	if err := preflightCheck(config); err != nil {
		log.Fatalln(err)
	}

	client := darksky.NewClient(
		config.GetString(CFG_API_KEY),
		config.GetDuration(CFG_CACHE_TTL),
	)

	lat, long, err := LookupCityCoords(config.GetString(CFG_CITY))
	if err != nil {
		log.Fatalln("Failed to lookup city:", err)
	}

	collectorOptions := DarkSkyCollectorOptions{
		City:         config.GetString(CFG_CITY),
		Lat:          forecast.Measurement(lat),
		LocationName: strings.ToLower(config.GetString(CFG_LOCATION_NAME)),
		Long:         forecast.Measurement(long),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(NewDarkSkyCollector(client, collectorOptions))

	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetInt(CFG_LISTEN_PORT)),
		Handler: m,
	}

	log.Println("Starting HTTP listener on", s.Addr)
	s.ListenAndServe()
}
