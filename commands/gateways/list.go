package gateways

import (
	"github.com/spf13/cobra"
	api "github.com/valyentdev/cli/api"
	tui "github.com/valyentdev/cli/tui"
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
