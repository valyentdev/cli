package fleets

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	fleetsCmd := &cobra.Command{
		Use: "fleets",
	}
	fleetsCmd.AddCommand(newCreateFleetCmd())
	fleetsCmd.AddCommand(newListFleetsCmd())

	return fleetsCmd
}
