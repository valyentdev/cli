package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/ravel/api"
)

func ListGateways(fleet api.Fleet) error {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Retrieve gateways from the API.
	gateways, err := client.GetGateways(fleet.Id)
	if err != nil {
		return err
	}

	// Compute the list of items.
	items := make([]list.Item, len(gateways))
	for idx, gateway := range gateways {
		// Retrieve fleet name
		items[idx] = ListItem{
			title:       gateway.Name,
			description: fmt.Sprintf("Fleet: %s | Target port: %d", fleet.Name, gateway.TargetPort),
		}
	}

	return List("List of Gateways", items)
}

func SelectGateway(fleet api.Fleet) (gateway *api.Gateway, err error) {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Retrieve gateways from the API.
	gateways, err := client.GetGateways(fleet.Id)
	if err != nil {
		return nil, err
	}

	// Compute the list of items.
	items := make([]list.Item, len(gateways))
	for idx, gateway := range gateways {
		// Retrieve fleet name

		items[idx] = FancySelectItem{
			title:       gateway.Name,
			description: fmt.Sprintf("Fleet: %s | Target port: %d", fleet.Name, gateway.TargetPort),
			value:       gateway.Id,
		}
	}

	gatewayID, err := FancySelect("Select Gateway", items)

	for _, gtw := range gateways {
		if gtw.Id == gatewayID {
			gateway = &gtw
		}
	}

	return
}
