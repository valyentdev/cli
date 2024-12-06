package api

import (
	"fmt"

	"github.com/valyentdev/cli/pkg/http"
	"github.com/valyentdev/ravel/api"
)

func GetMachines(fleetID string) ([]api.Machine, error) {
	machines := []api.Machine{}
	err := http.PerformRequest("GET", fmt.Sprintf("/v1/fleets/%s/machines", fleetID), nil, &machines)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve machines from the api: %v", err)
	}
	return machines, nil
}

func GetMachineEvents(fleetID, machineID string) ([]api.MachineEvent, error) {
	events := []api.MachineEvent{}
	err := http.PerformRequest("GET", fmt.Sprintf("/v1/fleets/%s/machines/%s/events", fleetID, machineID), nil, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events from the api: %v", err)
	}
	return events, nil
}
