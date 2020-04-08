package darksky

import (
	"fmt"
	forecast "github.com/mlbright/darksky/v2"
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

const (
	FORECAST_KEY_FORMAT  = "forecast-%s-%s-%s"
)

type Client struct {
  ApiKey string

	apiCache *cache.Cache
}

type Option func(*Client)

func NewClient(apiKey string, cacheTTL time.Duration, options ...Option) *Client {
  cli := &Client{
		ApiKey: apiKey,
		apiCache: cache.New(cacheTTL, 10*time.Minute),
	}

	for _, option := range options {
		option(cli)
	}

  return cli
}

func (c *Client) Get(lat, long, time string, units forecast.Units, lang forecast.Lang) (*forecast.Forecast, error) {
	cacheKey := fmt.Sprintf(FORECAST_KEY_FORMAT, lat, long, time)
	if data, found := c.apiCache.Get(cacheKey); found {
		return data.(*forecast.Forecast), nil
	}
	log.Printf("Fetching forecast for %s, %s @ %s", lat, long, time)

	data, err := forecast.Get(c.ApiKey, lat, long, time, units, lang)
	if err != nil {
		return nil, err
	}

	c.apiCache.Set(cacheKey, data, cache.DefaultExpiration)
	return data, nil
}
