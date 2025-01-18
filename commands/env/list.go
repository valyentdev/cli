package env

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newListEnvCmd() *cobra.Command {
	listEnvCmd := &cobra.Command{
		Use:   "list",
		Short: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}

			if err := runListEnvCmd(fleetID); err != nil {
				exit.WithError(err)
			}

			return nil
		},
	}
	listEnvCmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")

	return listEnvCmd
}

func runListEnvCmd(fleetID string) error {
	namespace, err := config.RetrieveNamespace()
	if err != nil {
		return err
	}

	if fleetID == "" {
		fleetID, err = tui.SelectFleet()
		if err != nil {
			return err
		}
	}

	client, err := http.NewClient()
	if err != nil {
		return err
	}

	envs, err := client.GetEnvironmentVariables(namespace, fleetID)
	if err != nil {
		return err
	}

	tui.ShowTable(
		[]table.Column{
			{Title: "Key", Width: 20},
			{Title: "Value", Width: 50},
		},
		envs,
		5,
	)

	return nil
}
