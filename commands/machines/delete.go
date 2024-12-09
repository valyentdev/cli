package machines

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/api"
	"github.com/valyentdev/cli/tui"
)

func newDeleteMachineCmd() *cobra.Command {
	logsCmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteMachineCmd()
		},
	}

	return logsCmd
}

func runDeleteMachineCmd() error {
	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	machineID, err := tui.SelectMachine(fleetID)
	if err != nil {
		return err
	}

	force := false
	if err := huh.NewConfirm().
		Title("Force delete?").
		Description("This will stop the machine if it is actually running.").
		Value(&force).
		Run(); err != nil {
		return err
	}

	if err := api.DeleteMachine(fleetID, machineID, force); err != nil {
		return err
	}

	fmt.Printf("Machine \"%s\" from fleet \"%s\" successfully deleted.\n", machineID, fleetID)

	return nil
}
