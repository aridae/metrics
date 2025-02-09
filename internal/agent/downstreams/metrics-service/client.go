package metricsservice

import (
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	client  *http.Client
	address string
}

func NewClient(address string, mws ...func(http.RoundTripper) http.RoundTripper) *Client {
	rt := http.DefaultTransport
	for _, mw := range mws {
		rt = mw(rt)
	}

	retryableHTTPClient := retryablehttp.NewClient()
	retryableHTTPClient.HTTPClient.Transport = rt

	return &Client{
		client:  retryableHTTPClient.StandardClient(),
		address: address,
	}
}

func doRequest(client *http.Client, method string, url string, body io.Reader, contentType string) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do http request: %w", err)
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}
