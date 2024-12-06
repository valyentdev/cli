package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/valyentdev/cli/internal/api"
	ravelAPI "github.com/valyentdev/ravel/api"
)

func ListGateways(fleets []ravelAPI.Fleet) error {
	// Retrieve gateways from the API.
	gateways, err := api.GetGateways()
	if err != nil {
		return err
	}

	// Compute the list of items.
	items := make([]list.Item, len(gateways))
	for idx, gateway := range gateways {
		// Retrieve fleet name
		fleetName := gateway.FleetId
		for _, fleet := range fleets {
			if fleet.Id == gateway.FleetId {
				fleetName = fleet.Name
			}
		}

		items[idx] = ListItem{
			title:       gateway.Name,
			description: fmt.Sprintf("Fleet: %s | Target port: %d", fleetName, gateway.TargetPort),
		}
	}

	return List("List of Gateways", items)
}
