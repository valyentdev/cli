package env

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	fleetsCmd := &cobra.Command{
		Use: "env",
	}
	fleetsCmd.AddCommand(newListEnvCmd())
	fleetsCmd.AddCommand(newSetEnvCmd())

	return fleetsCmd
}
