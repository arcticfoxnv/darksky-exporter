package darksky

import (
	"fmt"
	"github.com/arcticfoxnv/darksky-exporter/darksky/mock"
	forecast "github.com/shawntoffel/darksky"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClientForecast(t *testing.T) {
	s := mock.NewMockServer()
	defer s.Close()

	req := &forecast.ForecastRequest{
		Latitude:  40.7128,
		Longitude: -74.0059,
	}

	cli := NewClient("abc123", time.Minute, SetHTTPClient(s.Client()))
	cacheKey := fmt.Sprintf(FORECAST_KEY_FORMAT, req.Latitude, req.Longitude, req.Time)

	_, found := cli.apiCache.Get(cacheKey)
	assert.False(t, found)

	_, err := cli.Forecast(req)
	assert.Nil(t, err)

	_, found = cli.apiCache.Get(cacheKey)
	assert.True(t, found)

	_, err = cli.Forecast(req)
	assert.Nil(t, err)
}
