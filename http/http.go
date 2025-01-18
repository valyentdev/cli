package http

import (
	"fmt"
	"os"

	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/valyent.go"
)

// NewClient initializes a new Valyent HTTP API client instance,
// providing the API key and allowing to specify the API base url from an environment variable.
func NewClient() (*valyent.Client, error) {
	// Initialize a new Valyent HTTP API client instance.
	client := valyent.NewClient()

	// Customize the client's base URL field, if specified as an environment variable.
	if baseURL := os.Getenv("VALYENT_API_URL"); baseURL != "" {
		client.WithBaseURL(baseURL)
	}

	// Retrieve the API key
	apiKey, err := config.RetrieveAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve API key: %v", err)
	}

	// Set the retrieved API key as an authorization bearer token.
	client.WithBearerToken(apiKey)

	return client, nil
}
