package darksky

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	forecast "github.com/shawntoffel/darksky"
	"log"
	"net/http"
	"time"
)

const (
	FORECAST_KEY_FORMAT = "forecast-%g-%g-%d"
)

type Client struct {
	apiCache   *cache.Cache
	client     forecast.DarkSky
	httpClient *http.Client
}

type Option func(*Client)

func NewClient(apiKey string, cacheTTL time.Duration, options ...Option) *Client {
	cli := &Client{
		apiCache: cache.New(cacheTTL, 10*time.Minute),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, option := range options {
		option(cli)
	}

	cli.client = forecast.NewWithClient(apiKey, cli.httpClient)

	return cli
}

func (c *Client) Forecast(req *forecast.ForecastRequest) (forecast.ForecastResponse, error) {
	cacheKey := fmt.Sprintf(FORECAST_KEY_FORMAT, req.Latitude, req.Longitude, req.Time)
	if data, found := c.apiCache.Get(cacheKey); found {
		return data.(forecast.ForecastResponse), nil
	}
	log.Printf("Fetching forecast for %g, %g @ %d", req.Latitude, req.Longitude, req.Time)

	data, err := c.client.Forecast(*req)
	if err != nil {
		return forecast.ForecastResponse{}, err
	}

	c.apiCache.Set(cacheKey, data, cache.DefaultExpiration)
	return data, nil
}
