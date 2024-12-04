package commands

import (
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/commands/auth"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "valyent",
		Short: "A CLI tool to interact with Valyent's API.",
	}
	rootCmd.AddCommand(auth.NewCmd())
	rootCmd.AddCommand(newInitCmd())

	return rootCmd
}
