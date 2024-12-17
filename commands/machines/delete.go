package machines

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newDeleteMachineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"remove"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}

			machineID, err := cmd.Flags().GetString("machine")
			if err != nil {
				return err
			}

			if err := runDeleteMachineCmd(fleetID, machineID); err != nil {
				exit.WithError(err)
			}

			return nil
		},
	}
	cmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")
	cmd.Flags().StringP("machine", "m", "", "Machine's identifier (optional)")

	return cmd
}

func runDeleteMachineCmd(fleetID, machineID string) (err error) {
	// If the fleet if not specified by a flag, we let the user visually select it.
	if fleetID == "" {
		fleetID, err = tui.SelectFleet()
		if err != nil {
			return err
		}
	}

	if machineID == "" {
		machineID, err = tui.SelectMachine(fleetID)
		if err != nil {
			return err
		}
	}

	force := false
	if err := huh.NewConfirm().
		Title("Force delete?").
		Description("This will stop the machine if it is actually running.").
		Value(&force).
		Run(); err != nil {
		return err
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	if err := client.DeleteMachine(fleetID, machineID, force); err != nil {
		return err
	}

	fmt.Printf("Machine \"%s\" from fleet \"%s\" successfully deleted.\n", machineID, fleetID)

	return nil
}
