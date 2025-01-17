package gateways

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	gatewaysCmd := &cobra.Command{
		Use:     "gateways",
		Aliases: []string{"gateway", "gtw"},
	}
	gatewaysCmd.AddCommand(newCreateGatewayCmd())
	gatewaysCmd.AddCommand(newListGatewaysCmd())
	gatewaysCmd.AddCommand(newDeleteGatewayCmd())

	return gatewaysCmd
}
