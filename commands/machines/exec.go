package machines

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
	"github.com/valyentdev/ravel/api"
)

func newExecMachineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "exec",
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}

			machineID, err := cmd.Flags().GetString("machine")
			if err != nil {
				return err
			}

			if err := runExecMachineCmd(fleetID, machineID, args); err != nil {
				exit.WithError(err)
			}

			return nil
		},
	}
	cmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")
	cmd.Flags().StringP("machine", "m", "", "Machine's identifier (optional)")

	return cmd
}

func runExecMachineCmd(fleetID, machineID string, args []string) (err error) {
	// If the fleet is not specified by a flag, we let the user visually select it.
	if fleetID == "" {
		fleetID, err = tui.SelectFleet()
		if err != nil {
			return err
		}
	}

	// If the machine is not specified by a flag, we let the user visually select it.
	if machineID == "" {
		machineID, err = tui.SelectMachine(fleetID)
		if err != nil {
			return err
		}
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	if err := client.Exec(fleetID, machineID, api.ExecOptions{
		Cmd: args,
	}); err != nil {
		return err
	}

	fmt.Printf("Command successfully executed")

	return nil
}
