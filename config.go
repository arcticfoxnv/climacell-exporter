package main

import (
	"github.com/spf13/viper"
	"os"
)

const (
	CFG_LISTEN_PORT = "listen_port"
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

	// Read config
	if err := v.ReadInConfig(); err != nil {
		return v, err
	}

	return v, nil
}

func preflightCheck(v *viper.Viper) error {
	// Check if required values have been set, return error if not

	// Return nil if all required values are set
	return nil
}
