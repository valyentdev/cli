package api

import (
	"fmt"
	stdHTTP "net/http"

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

func DeleteGateway(gatewayID string) error {
	err := http.PerformRequest(
		stdHTTP.MethodDelete,
		"/v1/gateways/"+gatewayID,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to delete gateway: %v", err)
	}
	return err
}
