package webhook

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	url        string
	httpClient *http.Client
	headers    map[string]string
}

func (c *Client) build(ctx context.Context, request any) (*http.Request, error) {
	if request == nil {
		return http.NewRequestWithContext(ctx, http.MethodPost, c.url, nil)
	}

	var reqBytes []byte
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.url,
		bytes.NewBuffer(reqBytes),
	)
}

func (c *Client) Send(ctx context.Context, payload any) error {
	req, err := c.build(ctx, payload)
	if err != nil {
		return &RequestError{
			HTTPStatusCode: http.StatusInternalServerError,
			Err:            fmt.Errorf("build request with error: %s", err.Error()),
		}
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Add headers to request
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return &RequestError{
			HTTPStatusCode: http.StatusInternalServerError,
			Err:            fmt.Errorf("request failed with error: %s", err.Error()),
		}
	}
	defer res.Body.Close()

	if isFailureStatusCode(res) {
		return &RequestError{
			HTTPStatusCode: res.StatusCode,
			Err:            fmt.Errorf("request failed with status code: %d", res.StatusCode),
		}
	}

	return nil
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func NewClient(s string, insecure bool, headers map[string]string) *Client {
	if s == "" {
		return nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil
	}

	client := http.DefaultClient
	client.Timeout = 5 * time.Second

	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		url:        s,
		httpClient: client,
		headers:    headers,
	}
}
