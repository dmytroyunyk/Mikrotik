package mikrotik

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	baseURL  string
	username string
	password string
	retries  int
	http     *http.Client
	mu       sync.RWMutex
	healthy  bool
}

func NewClient(host, username, password string, timeout time.Duration, retries int) *Client {
	return &Client{
		baseURL:  fmt.Sprintf("http://%s/rest", host),
		username: username,
		password: password,
		retries:  retries,
		healthy:  true,
		http: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	url := c.baseURL + path

	var lastErr error
	for attempt := range c.retries {
		if attempt > 0 {
			wait := time.Duration(attempt) * 500 * time.Millisecond
			slog.Warn("retry the request", "path", path, "attempt", attempt+1, "wait", wait)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("error creating query: %w", err)
		}
		req.SetBasicAuth(c.username, c.password)
		req.Header.Set("Accept", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("connection error: %w", err)
			c.setHealth(false)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication error: check login/password")
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("API return %d: %s", resp.StatusCode, body)
		}

		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("return answer: %w", err)
		}

		c.setHealth(true)
		return nil
	}

	return fmt.Errorf("all %d attempts exhausted, last error: %w", c.retries, lastErr)
}

func (c *Client) setHealth(ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.healthy = ok
}

func (c *Client) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.healthy
}

func (c *Client) Ping(ctx context.Context) error {
	var result map[string]any
	return c.Get(ctx, "/system/identity", &result)
}
