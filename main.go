package main

import (
	"darksky-exporter/darksky"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	DefaultCacheTTL   = 5 * time.Minute
	DefaultListenPort = 8080
)

func getEnvWithDefault(key, defaultValue string) string {
	if value, present := os.LookupEnv(key); present {
		return value
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value, present := os.LookupEnv(key); present {
		v, _ := strconv.Atoi(value)
		return v
	}
	return defaultValue
}

func preflightCheck(config *Config) error {
	if config.ApiKey == "" {
		return errors.New("Cannot start exporter, api key is missing")
	}

	if config.City == "" {
		return errors.New("Cannot start exporter, city not set")
	}

	return nil
}

func main() {
	log.Printf("darksky-exporter v%s-%s", Version, Commit)
	log.Println("Powered by Dark Sky - https://darksky.net/poweredby/")
	cfgFilename := getEnvWithDefault("DARKSKY_CONFIG_FILE", "darksky.toml")
	config, err := LoadConfig(cfgFilename)
	if err != nil {
		log.Println("Failed to load config file:", err)
		config = &Config{}
	}
	setConfigDefaults(config)

	config.ApiKey = getEnvWithDefault("DARKSKY_API_KEY", config.ApiKey)
	config.City = getEnvWithDefault("DARKSKY_CITY", config.City)
	listenPort := getEnvIntWithDefault("DARKSKY_LISTEN", DefaultListenPort)

	if err := preflightCheck(config); err != nil {
		log.Fatalln(err)
	}

	client := darksky.NewClient(
		config.ApiKey,
		config.CacheTTL,
	)

	lat, long, err := LookupCityCoords(config.City)
	if err != nil {
		log.Fatalln("Failed to lookup city:", err)
	}

	collectorOptions := DarkSkyCollectorOptions{
		City: FormatCityName(config.City),
		Lat:  fmt.Sprintf("%f", lat),
		Long: fmt.Sprintf("%f", long),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(NewDarkSkyCollector(client, collectorOptions))

	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: m,
	}

	log.Println("Starting HTTP listener on", s.Addr)
	s.ListenAndServe()
}
