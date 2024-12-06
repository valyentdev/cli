package gateways

import (
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/api"
	"github.com/valyentdev/cli/internal/tui"
)

func newListGatewaysCmd() *cobra.Command {
	listGatewaysCmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListGatewaysCmd()
		},
	}

	return listGatewaysCmd
}

func runListGatewaysCmd() error {
	// Retrieve fleets from the API.
	fleets, err := api.GetFleets()
	if err != nil {
		return err
	}

	// List gateways.
	if err := tui.ListGateways(fleets); err != nil {
		return err
	}

	return nil
}
