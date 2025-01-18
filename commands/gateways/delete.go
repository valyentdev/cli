package gateways

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newDeleteGatewayCmd() *cobra.Command {
	deleteGatewayCmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"remove"},
		RunE: func(cmd *cobra.Command, args []string) error {
			gatewayID, err := cmd.Flags().GetString("gateway")
			if err != nil {
				return err
			}
			confirmed, err := cmd.Flags().GetBool("confirmed")
			if err != nil {
				return err
			}
			if err := runDeleteGatewayCmd(gatewayID, confirmed); err != nil {
				exit.WithError(err)
			}
			return nil
		},
	}

	deleteGatewayCmd.Flags().StringP("gateway", "g", "", "Gateway's identifier (optional)")
	deleteGatewayCmd.Flags().BoolP("confirmed", "c", false, "Confirm gateway deletion (optional)")

	return deleteGatewayCmd
}

func runDeleteGatewayCmd(gatewayID string, confirmed bool) (err error) {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	fleet, err := client.GetFleet(fleetID)
	if err != nil {
		return err
	}

	// If the gateway is not specified by the user as a command flag,
	// we ask for it with a nice TUI.
	if gatewayID == "" {
		gtw, err := tui.SelectGateway(fleet)
		if err != nil {
			return fmt.Errorf("failed to select gateway: %v", err)
		}
		gatewayID = gtw.Id
	}

	// Ask for deletion confirmation,
	// if not already specified as a flag to the command.
	if !confirmed {
		err = huh.
			NewConfirm().
			Title("Are you sure you want to delete this gateway?").
			Value(&confirmed).
			Run()
		if err != nil {
			return err
		}
	}

	// If the deletion is still not confirmed by the user,
	// we skip the API call...
	if !confirmed {
		fmt.Println("Cancelling gateway deletion...")
		return
	}

	// Now, we can safely try to delete the gateway.
	if err := client.DeleteGateway(fleetID, gatewayID); err != nil {
		return err
	}

	fmt.Println("Gateway successfully deleted!")

	return
}
