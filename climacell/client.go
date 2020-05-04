package climacell

import (
	"fmt"
	api "github.com/arcticfoxnv/climacell-go"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"
)

const (
	FORECAST_KEY_FORMAT = "forecast-%g-%g-%s"
)

type Client struct {
	apiCache   *cache.Cache
	client     *api.Client
	httpClient *http.Client
}

type Option func(*Client)

func SetHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

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

	cli.client = api.NewClient(apiKey, api.SetHTTPClient(cli.httpClient))

	return cli
}

func (c *Client) SetUserAgent(ua string) {
	c.client.UserAgent = ua
}

func (c *Client) RealtimeWeather(req *api.RealtimeRequest) (*api.RealtimeResponse, error) {
	cacheKey := fmt.Sprintf(FORECAST_KEY_FORMAT, req.Latitude, req.Longitude, "now")
	if data, found := c.apiCache.Get(cacheKey); found {
		return data.(*api.RealtimeResponse), nil
	}
	log.Printf("Fetching forecast for %g, %g", req.Latitude, req.Longitude)

	data, err := c.client.RealtimeWeather(req)
	if err != nil {
		return nil, err
	}

	c.apiCache.Set(cacheKey, data, cache.DefaultExpiration)
	return data, nil
}
