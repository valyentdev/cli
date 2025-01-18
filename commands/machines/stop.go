package machines

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newStopMachineCmd() *cobra.Command {
	logsCmd := &cobra.Command{
		Use: "stop",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runStopMachineCmd(); err != nil {
				exit.WithError(err)
			}
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

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking to stop this machine.
	if err := client.StopMachine(fleetID, machineID); err != nil {
		return fmt.Errorf("the API call asking to stop machine failed: %v", err)
	}

	fmt.Println("Machine stopped!")

	return nil
}
