package main

import (
	"bytes"
	"github.com/arcticfoxnv/climacell-exporter/climacell"
	"github.com/arcticfoxnv/climacell-exporter/climacell/mock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {

	// Read expected results from file
	data, err := ioutil.ReadFile("testdata/metrics.txt")
	if err != nil {
		t.Fail()
	}

	expected := bytes.NewReader(data)

	// Setup collector
	s := mock.NewMockServer()
	defer s.Close()

	client := climacell.NewClient("abc123", time.Minute, climacell.SetHTTPClient(s.Client()))
	c := NewCollector(client, CollectorOptions{
		City:         "New York, NY",
		Latitude:          40.7128,
		LocationName: "test",
		Longitude:         -74.0059,
		EnableWeatherDataLayer: true,
	})

	// Test collector and check for errors
	err = testutil.CollectAndCompare(c, expected)
	assert.Nil(t, err)
}
