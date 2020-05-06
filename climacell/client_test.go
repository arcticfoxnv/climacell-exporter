package climacell

import (
	"fmt"
	"github.com/arcticfoxnv/climacell-exporter/climacell/mock"
	api "github.com/arcticfoxnv/climacell-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClientRealtimeWeather(t *testing.T) {
	s := mock.NewMockServer()
	defer s.Close()

	req := &api.RealtimeRequest{
		Latitude:  40.7128,
		Longitude: -74.0059,
		Fields: api.DataFieldList{
			api.EpaAqi,
			api.Dewpoint,
			api.FeelsLike,
		},
	}

	cli := NewClient("abc123", time.Minute, SetHTTPClient(s.Client()))
	cacheKey := fmt.Sprintf(FORECAST_KEY_FORMAT, req.Latitude, req.Longitude, "now")

	_, found := cli.apiCache.Get(cacheKey)
	assert.False(t, found)

	_, err := cli.RealtimeWeather(req)
	assert.Nil(t, err)

	_, found = cli.apiCache.Get(cacheKey)
	assert.True(t, found)

	_, err = cli.RealtimeWeather(req)
	assert.Nil(t, err)
}
