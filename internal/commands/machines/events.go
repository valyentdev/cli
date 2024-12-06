package machines

import (
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/tui"
)

func newListMachineEventsCmd() *cobra.Command {
	listGatewaysCmd := &cobra.Command{
		Use: "events",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListMachineEventsCmd()
		},
	}

	return listGatewaysCmd
}

func runListMachineEventsCmd() error {
	// Select the fleet from which we want to list machines.
	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	machineID, err := tui.SelectMachine(fleetID)
	if err != nil {
		return err
	}

	if err := tui.ListMachineEvents(fleetID, machineID); err != nil {
		return err
	}

	return nil
}
