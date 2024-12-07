package api

import (
	"fmt"

	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/ravel/api"
)

func GetGateways() ([]api.Gateway, error) {
	gateways := []api.Gateway{}
	err := http.PerformRequest("GET", "/v1/gateways", nil, &gateways)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gateways from the api: %v", err)
	}
	return gateways, nil
}
