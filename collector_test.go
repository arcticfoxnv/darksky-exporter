package main

import (
	"bytes"
	"github.com/arcticfoxnv/darksky-exporter/darksky"
	"github.com/arcticfoxnv/darksky-exporter/darksky/mock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestDarkSkyCollector(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/metrics.txt")
	if err != nil {
		t.Fail()
	}

	expected := bytes.NewReader(data)

	s := mock.NewMockServer()
	defer s.Close()

	client := darksky.NewClient("abc123", time.Minute, darksky.SetHTTPClient(s.Client()))
	c := NewDarkSkyCollector(client, DarkSkyCollectorOptions{
		City:         "New York, NY",
		Lat:          40.7128,
		LocationName: "test",
		Long:         -74.0059,
	})

	err = testutil.CollectAndCompare(c, expected)
	assert.Nil(t, err)
}
