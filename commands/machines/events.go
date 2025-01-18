package machines

import (
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newListMachineEventsCmd() *cobra.Command {
	listGatewaysCmd := &cobra.Command{
		Use: "events",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runListMachineEventsCmd(); err != nil {
				exit.WithError(err)
			}
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
