package metricsservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) UpdateMetricsBatch(_ context.Context, metrics []Metric) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(metrics); err != nil {
		return fmt.Errorf("failed to encode json-serializable struct: %w", err)
	}

	serverURL, _ := url.JoinPath("http://"+c.address, "/update")

	if err := doRequest(c.client, http.MethodPost, serverURL, body, "application/json"); err != nil {
		return fmt.Errorf("failed to do http call: %w", err)
	}

	return nil
}
