package config

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"os"
)

type (
	Config struct {
		UrlConfig UrlConfig `json:"url_config"`
		Logging   Logging   `json:"logging"`
	}
	RateLimitConfig struct {
		IntervalSeconds int `json:"interval_seconds"`
		MaxRequests     int `json:"max_requests"`
	}
	UrlConfig struct {
		RateLimit RateLimitConfig `json:"rate_limit"`
		URLs      []string        `json:"urls"`
	}
	Logging struct {
		LogServiceName string        `json:"log_service_name"`
		LogLevel       zerolog.Level `json:"log_level"`
	}
)

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
