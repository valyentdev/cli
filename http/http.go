package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/pkg/env"
)

// PerformRequest sends an HTTP request with an optional JSON payload and unmarshals the JSON response if responseTarget is non-nil.
func PerformRequest(method, path string, payload any, responseTarget any) error {
	// Construct the URL
	baseURL := env.GetVar("VALYENT_API_URL", "https://console.valyent.cloud")
	url := baseURL + path

	// Prepare request body if payload is provided
	var body io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload to JSON: %w", err)
		}
		body = bytes.NewReader(payloadBytes)
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add authorization header if API key is available
	apiKey, _ := config.RetrieveAPIKey()
	if apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+apiKey)
	}

	// Set appropriate headers
	if payload != nil {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx HTTP response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %s", resp.Status)
	}

	// Read and process the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// If responseTarget is non-nil, unmarshal the JSON response into it
	if responseTarget != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, responseTarget); err != nil {
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}

	return nil
}
