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

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking to start this machine.
	if err := client.StartMachine(fleetID, machineID); err != nil {
		return fmt.Errorf("the API call asking to start machine failed: %v", err)
	}

	fmt.Println("Machine started!")

	return nil
}
