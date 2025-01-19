package fleets

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	fleetsCmd := &cobra.Command{
		Use:     "fleets",
		Aliases: []string{"fleet"},
	}
	fleetsCmd.AddCommand(newCreateFleetCmd())
	fleetsCmd.AddCommand(newListFleetsCmd())
	fleetsCmd.AddCommand(newDeleteFleetCmd())

	return fleetsCmd
}
