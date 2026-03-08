package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string

	rateLimiter <-chan time.Time
	
	maxRetries int
	retryDelay time.Duration
}

func New() *Client {
	ticker := time.NewTicker(3 * time.Second)

	transport := &http.Transport{
		MaxIdleConns: 1,
		MaxIdleConnsPerHost: 1,
		IdleConnTimeout: 90 * time.Second,
	}

	return &Client {
		baseURL: "http://export.arxiv.org/api/query",
		userAgent: "arxiv-mcp-server/1.0.0",

		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: transport,
		},

		rateLimiter: ticker.C,

		maxRetries: 3,
		retryDelay: 3 * time.Second,
	}
}

func (c *Client) do(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	var respErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case <- c.rateLimiter:
		}

		resp, respErr := c.httpClient.Do(req)
		if respErr == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		defer resp.Body.Close()
		
		if attempt == c.maxRetries {
			break
		}

		backof := time.Duration(1<<attempt) * c.retryDelay
		time.Sleep(backof)
	}

	return resp, respErr
}

func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, nil, nil)
}

func ReadBody(resp *http.Response) ([]byte, error){
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}