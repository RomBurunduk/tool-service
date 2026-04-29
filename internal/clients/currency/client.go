package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"tool-service/internal/config"
)

type Client struct {
	httpClient *http.Client
	url        string
	apiKey     string
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: cfg.CurrencyHTTPTimeout},
		url:        cfg.CurrencyAPIURL,
		apiKey:     cfg.CurrencyAPIKey,
	}
}

func (c *Client) FetchLatest(ctx context.Context) (LatestAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return LatestAPIResponse{}, err
	}
	req.Header.Set("apikey", strings.TrimSpace(c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LatestAPIResponse{}, fmt.Errorf("currency http: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LatestAPIResponse{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return LatestAPIResponse{}, fmt.Errorf("currency status %d: %s", resp.StatusCode, truncateStr(string(body), 512))
	}
	var out LatestAPIResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return LatestAPIResponse{}, fmt.Errorf("currency json: %w", err)
	}
	return out, nil
}

func truncateStr(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
