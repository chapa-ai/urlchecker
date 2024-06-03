package checks

import (
	"context"
	"github.com/chapa-ai/urlchecker/config"
	"github.com/rs/zerolog"
	"os"
	"reflect"
	"testing"
	"time"
)

// / testing the URL status check function
func TestCheckStatusCodeAndText(t *testing.T) {
	ctx := context.Background()
	log := zerolog.New(os.Stderr)

	testCases := []struct {
		url           string
		expectedText  string
		expectedError bool
	}{
		{"https://www.google.com", "ok", false},
		{"https://nonexistent.url", "fail", true},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			result, err := checkStatusCodeAndText(ctx, log, tc.url)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
			if result != nil && result.text != tc.expectedText {
				t.Errorf("expected text: %v, got: %v", tc.expectedText, result.text)
			}
		})
	}
}

// / testing the configuration file reading function
func TestLoadConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `{
      "url_config": {
        "rate_limit": {
          "interval_seconds": 10,
          "max_requests": 10
        },
        "urls": [
          "https://vk.com",
          "https://ya.ru"
        ]
      },
      "logging": {
        "log_service_name": "urlchecker",
        "log_level": "info"
      }
    }`

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.UrlConfig.RateLimit.IntervalSeconds != 10 || cfg.UrlConfig.RateLimit.MaxRequests != 10 {
		t.Errorf("unexpected rate limit config: %+v", cfg.UrlConfig.RateLimit)
	}

	expectedURLs := []string{"https://vk.com", "https://ya.ru"}
	if !reflect.DeepEqual(cfg.UrlConfig.URLs, expectedURLs) {
		t.Errorf("expected URLs: %v, got: %v", expectedURLs, cfg.UrlConfig.URLs)
	}

	if cfg.Logging.LogServiceName != "urlchecker" || cfg.Logging.LogLevel != zerolog.InfoLevel {
		t.Errorf("unexpected logging config: %+v", cfg.Logging)
	}
}

// / testing processing a list of URLs
func TestDoChecksWithInterval(t *testing.T) {
	log := zerolog.New(os.Stderr)
	cfg := config.Config{
		UrlConfig: config.UrlConfig{
			RateLimit: config.RateLimitConfig{
				IntervalSeconds: 1,
				MaxRequests:     2,
			},
			URLs: []string{"https://www.google.com", "https://www.bing.com"},
		},
		Logging: config.Logging{
			LogServiceName: "test_service",
			LogLevel:       zerolog.InfoLevel,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error)

	go func() {
		err := DoChecksWithInterval(log, cfg)
		if err != nil {
			done <- err
			return
		}
		done <- nil
	}()

	select {
	case <-ctx.Done():
		t.Log("test finished due to timeout")
	case err := <-done:
		if err != nil {
			t.Fatalf("failed DoChecksWithInterval: %v", err)
		} else {
			t.Log("DoChecksWithInterval completed successfully")
		}
	}
}
