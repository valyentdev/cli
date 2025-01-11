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

			return runListEnvCmd(fleetID)
		},
	}
	listEnvCmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")

	return listEnvCmd
}

func runListEnvCmd(fleetID string) (err error) {
	namespace, err := config.RetrieveNamespace()
	if err != nil {
		exit.WithError(err)
	}

	if fleetID == "" {
		fleetID, err = tui.SelectFleet()
		if err != nil {
			exit.WithError(err)
		}
	}

	client, err := http.NewClient()
	if err != nil {
		exit.WithError(err)
	}

	envs, err := client.GetEnvironmentVariables(namespace, fleetID)
	if err != nil {
		exit.WithError(err)
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
