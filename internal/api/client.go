package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	http      *retryablehttp.Client
	cb        *CircuitBreaker
	token     string
	baseURL   string
	workspace string
	project   string
}

func NewClient(token, baseURL, workspace, project string, cfg *RetryConfig) *Client {
	return &Client{
		http:      newRetryableHTTP(cfg),
		cb:        NewCircuitBreaker(cfg.CircuitBreakerLimit),
		token:     token,
		baseURL:   strings.TrimRight(baseURL, "/"),
		workspace: workspace,
		project:   project,
	}
}

func (c *Client) Do(method, path string, body interface{}, query url.Values) ([]byte, int, error) {
	if c.cb.IsOpen() {
		return nil, 0, ErrCircuitBreakerOpen
	}

	u := fmt.Sprintf("%s/api/v1/workspaces/%s/projects/%s/%s", c.baseURL, c.workspace, c.project, path)
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := retryablehttp.NewRequest(method, u, reqBody)
	if err != nil {
		c.cb.RecordFailure()
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-Key", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "plane-cli/0.1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		c.cb.RecordFailure()
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.cb.RecordFailure()
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.cb.RecordSuccess()
		return respBody, resp.StatusCode, nil
	}

	c.cb.RecordFailure()
	return respBody, resp.StatusCode, fmt.Errorf("API error: %s", resp.Status)
}
