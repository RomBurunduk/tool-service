package wordstat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tool-service/internal/config"
	"tool-service/internal/model"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	path       string
	apiKey     string
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: cfg.WordStatHTTPTimeout},
		baseURL:    strings.TrimRight(cfg.WordStatBaseURL, "/"),
		path:       cfg.WordStatPath,
		apiKey:     cfg.WordStatAPIKey,
	}
}

type dynamicsRequest struct {
	Phrase   string `json:"phrase"`
	Period   string `json:"period"`
	FromDate string `json:"fromDate"`
	FolderID string `json:"folderId"`
}

func (c *Client) Dynamics(ctx context.Context, phrase, folderID string, fromDate time.Time) (model.WordStatResult, error) {
	reqBody := dynamicsRequest{
		Phrase:   phrase,
		Period:   "PERIOD_WEEKLY",
		FromDate: fromDate.UTC().Format(time.RFC3339),
		FolderID: folderID,
	}
	raw, err := json.Marshal(reqBody)
	if err != nil {
		return model.WordStatResult{}, err
	}

	u, err := url.Parse(c.baseURL + c.path)
	if err != nil {
		return model.WordStatResult{}, fmt.Errorf("wordstat url: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(raw))
	if err != nil {
		return model.WordStatResult{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Authorization", "Api-Key "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return model.WordStatResult{}, fmt.Errorf("wordstat http: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.WordStatResult{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return model.WordStatResult{}, fmt.Errorf("wordstat status %d: %s", resp.StatusCode, truncate(body, 512))
	}
	result := model.WordStatResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return model.WordStatResult{}, err
	}
	return result, err
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
