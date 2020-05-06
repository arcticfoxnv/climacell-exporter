package main

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	CFG_API_KEY       = "api_key"
	CFG_CACHE_TTL     = "cache_ttl"
	CFG_CITY          = "city"
	CFG_LISTEN_PORT   = "listen_port"
	CFG_LOCATION_NAME = "location_name"
)

var (
	ErrApiKeyMissing       = errors.New("Cannot start exporter, api key is missing")
	ErrCityMissing         = errors.New("Cannot start exporter, city not set")
	ErrLocationNameMissing = errors.New("Cannot start exporter, location name not set")
)

func loadConfig() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper
	v.SetConfigName("climacell")
	v.SetConfigType("toml")
	v.AddConfigPath("/etc")
	v.AddConfigPath(".")
	v.SetEnvPrefix("climacell")
	v.AutomaticEnv()

	if path, present := os.LookupEnv("CLIMACELL_CONFIG_FILE"); present {
		v.SetConfigFile(path)
	}

	// Configure defaults
	v.SetDefault(CFG_LISTEN_PORT, 8080)
	v.SetDefault(CFG_CACHE_TTL, 5*time.Minute)

	// Read config
	if err := v.ReadInConfig(); err != nil {
		return v, err
	}

	return v, nil
}

func preflightCheck(v *viper.Viper) error {
	// Check if required values have been set, return error if not
	if v.GetString(CFG_API_KEY) == "" {
		return ErrApiKeyMissing
	}

	if v.GetString(CFG_CITY) == "" {
		return ErrCityMissing
	}

	if v.GetString(CFG_LOCATION_NAME) == "" {
		return ErrLocationNameMissing
	}

	// Return nil if all required values are set
	return nil
}
