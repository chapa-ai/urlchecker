package checks

import (
	"context"
	"github.com/chapa-ai/urlchecker/config"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
	"time"
)

const (
	HttpMethodGet = "GET"
)

type Response struct {
	url        string
	statusCode int
	text       string
}

func checkStatusCodeAndText(ctx context.Context, log zerolog.Logger, url string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, HttpMethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create request: %s", url)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get response: %s", url)
		return nil, err
	}
	defer resp.Body.Close()

	text := "ok"
	if resp.StatusCode != 200 {
		text = "fail"
	}

	result := &Response{
		url:        url,
		statusCode: resp.StatusCode,
		text:       text,
	}

	return result, nil
}

func DoChecksWithInterval(log zerolog.Logger, urls config.Config) error {
	ticker := time.NewTicker(time.Second * time.Duration(urls.UrlConfig.RateLimit.IntervalSeconds))
	defer ticker.Stop()
	sem := make(chan struct{}, urls.UrlConfig.RateLimit.MaxRequests)
	var wg sync.WaitGroup

	for range ticker.C {
		for _, url := range urls.UrlConfig.URLs {
			sem <- struct{}{}
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				defer func() { <-sem }()

				log.Info().Msgf("checking url: %s", url)
				results, err := checkStatusCodeAndText(context.Background(), log, url)
				if err != nil {
					log.Error().Err(err).Msgf("failed CheckStatusCodeAndText: %s", err)
					return
				}
				log.Info().Msgf("%s: %s - %d\n", results.url, results.text, results.statusCode)

				log.Info().Msgf("finished checking url: %s", url)
			}(url)
		}
		wg.Wait()
	}
	return nil
}
