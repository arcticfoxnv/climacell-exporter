package main

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCollector(t *testing.T) {

	// Read expected results from file
	data, err := ioutil.ReadFile("testdata/metrics.txt")
	if err != nil {
		t.Fail()
	}

	expected := bytes.NewReader(data)

	// Setup collector
	c := NewCollector(CollectorOptions{})

	// Test collector and check for errors
	err = testutil.CollectAndCompare(c, expected)
	assert.Nil(t, err)
}
