package darksky

import (
	"context"
	"crypto/tls"
	"fmt"
  forecast "github.com/mlbright/darksky/v2"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
  "time"
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

func TestClientGet(t *testing.T) {
  lat := "0"
  long := "0"
  when := "now"
  units := forecast.AUTO
  lang := forecast.English


	cli := NewClient("abc123", time.Minute)
  cacheKey := fmt.Sprintf(FORECAST_KEY_FORMAT, lat, long, when)

  _, found := cli.apiCache.Get(cacheKey)
	assert.False(t, found)

	_, err := cli.Get(lat, long, when, units, lang)
  assert.Nil(t, err)

  _, found = cli.apiCache.Get(cacheKey)
	assert.True(t, found)

  _, err = cli.Get(lat, long, when, units, lang)
  assert.Nil(t, err)
}
