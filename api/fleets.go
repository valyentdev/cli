package api

import (
	"fmt"
	stdHTTP "net/http"

	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/ravel/api"
)

type CreateFleetOptions struct {
	Name string `json:"name"`
}

func CreateFleet(name string) error {
	err := http.PerformRequest(
		"POST",
		"/v1/fleets",
		CreateFleetOptions{
			Name: name,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create fleet: %v", err)
	}

	return nil
}

func GetFleets() (fleets []api.Fleet, err error) {
	// Fetch existing fleets matching the user's namespace.
	fleets = []api.Fleet{}
	err = http.PerformRequest(stdHTTP.MethodGet, "/v1/fleets", nil, &fleets)
	if err != nil {
		return fleets, fmt.Errorf("failed to retrieve fleets: %v", err)
	}

	return
}

func DeleteFleet(fleetID string) error {
	err := http.PerformRequest(stdHTTP.MethodDelete, "/v1/fleets/"+fleetID, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete fleet: %v", err)
	}

	return nil
}
