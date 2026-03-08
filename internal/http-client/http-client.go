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
}

func New() *Client {
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

	return c.httpClient.Do(req)
}

func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, nil, nil)
}

func ReadBody(resp *http.Response) ([]byte, error){
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}