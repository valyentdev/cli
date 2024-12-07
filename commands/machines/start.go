package machines

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/tui"
)

func newStartMachineCmd() *cobra.Command {
	logsCmd := &cobra.Command{
		Use: "start",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStartMachineCmd()
		},
	}

	return logsCmd
}

func runStartMachineCmd() error {
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
		Title("Are you sure you want to start this machine?").
		Value(&confirmed).
		Run()
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	err = http.PerformRequest("POST", fmt.Sprintf("/fleets/%s/machines/%s/start", fleetID, machineID), nil, nil)
	if err != nil {
		return err
	}

	fmt.Println("Machine started!")

	return nil
}
