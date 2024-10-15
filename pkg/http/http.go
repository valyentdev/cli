package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/valyentdev/cli/pkg/env"
)

// PerformRequest sends an HTTP request and unmarshalls the JSON response.
func PerformRequest(method, path string, v any) error {
	// Create HTTP client
	client := &http.Client{}

	// Create new request
	baseURL := env.GetVar("VALYENT_API_URL", "https://console.valyent.dev")
	url := baseURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	// Unmarshal JSON into provided interface
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return nil
}
