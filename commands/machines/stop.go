package machines

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/tui"
)

func newStopMachineCmd() *cobra.Command {
	logsCmd := &cobra.Command{
		Use: "stop",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStopMachineCmd()
		},
	}

	return logsCmd
}

func runStopMachineCmd() error {
	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	machineID, err := tui.SelectMachine(fleetID)
	if err != nil {
		return err
	}

	confirmed := false
	err = huh.NewConfirm().
		Title("Are you sure you want to stop this machine?").
		Value(&confirmed).
		Run()
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	err = http.PerformRequest("POST", fmt.Sprintf("/fleets/%s/machines/%s/stop", fleetID, machineID), nil, nil)
	if err != nil {
		return err
	}

	fmt.Println("Machine stopped!")

	return nil
}
