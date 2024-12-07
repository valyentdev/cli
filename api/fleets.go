package api

import (
	"fmt"
	stdHTTP "net/http"

	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/ravel/api"
)

func GetFleets() (fleets []api.Fleet, err error) {
	// Fetch existing fleets matching the user's namespace.
	fleets = []api.Fleet{}
	err = http.PerformRequest(stdHTTP.MethodGet, "/v1/fleets", nil, &fleets)
	if err != nil {
		return fleets, fmt.Errorf("failed to retrieve fleets: %v", err)
	}

	return
}
