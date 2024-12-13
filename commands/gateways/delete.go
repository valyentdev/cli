package gateways

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/tui"
)

func newDeleteGatewayCmd() *cobra.Command {
	deleteGatewayCmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			gatewayID, err := cmd.Flags().GetString("gateway")
			if err != nil {
				return err
			}
			confirmed, err := cmd.Flags().GetBool("confirmed")
			if err != nil {
				return err
			}
			return runDeleteGatewayCmd(gatewayID, confirmed)
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

	// We retrieve fleets from the API,
	// from the currently authenticated namespace.
	fleets, err := client.GetFleets()
	if err != nil {
		return err
	}

	// If the gateway is not specified by the user as a command flag,
	// we ask for it with a nice TUI.
	if gatewayID == "" {
		gatewayID, err = tui.SelectGateway(fleets)
		if err != nil {
			return err
		}
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
	if err := client.DeleteGateway(gatewayID); err != nil {
		return err
	}

	fmt.Println("Gateway successfully deleted!")

	return
}
