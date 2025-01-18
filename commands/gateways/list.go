package gateways

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/valyentdev/cli/http"
	tui "github.com/valyentdev/cli/tui"
)

func newListGatewaysCmd() *cobra.Command {
	listGatewaysCmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListGatewaysCmd()
		},
	}

	return listGatewaysCmd
}

func runListGatewaysCmd() error {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	fleet, err := client.GetFleet(fleetID)
	if err != nil {
		return err
	}

	// List gateways.
	if err := tui.ListGateways(fleet); err != nil {
		return err
	}

	return nil
}
