package commands

import (
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/commands/auth"
	"github.com/valyentdev/cli/internal/commands/fleets"
	"github.com/valyentdev/cli/internal/commands/gateways"
	"github.com/valyentdev/cli/internal/commands/machines"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "valyent",
		Short: "A CLI tool to interact with Valyent's API.",
	}
	rootCmd.AddCommand(auth.NewCmd())
	rootCmd.AddCommand(fleets.NewCmd())
	rootCmd.AddCommand(gateways.NewCmd())
	rootCmd.AddCommand(machines.NewCmd())
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newDeployCmd())

	return rootCmd
}
