package machines

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/api"
	"github.com/valyentdev/cli/tui"
)

func newListMachinesCmd() *cobra.Command {
	listGatewaysCmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListMachinesCmd()
		},
	}

	return listGatewaysCmd
}

func runListMachinesCmd() error {
	// Retrieve fleets from the API.
	fleets, err := api.GetFleets()
	if err != nil {
		return fmt.Errorf("failed to retrieve fleets from the api: %v", err)
	}

	// Select the fleet from which we want to list machines.
	fleetID, err := tui.SelectFleetWithFleets(fleets)
	if err != nil {
		return fmt.Errorf("failed to select fleet: %v", err)
	}

	// List machines.
	if err := tui.ListMachines(fleets, fleetID); err != nil {
		return fmt.Errorf("failed to list machines: %v", err)
	}

	return nil
}
