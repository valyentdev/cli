package env

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newSetEnvCmd() *cobra.Command {
	setEnvCmd := &cobra.Command{
		Use:     "set <KEY>=<VALUE>",
		Short:   "Set environment",
		Example: "valyent env set DATABASE_URL=postgresql://username:password@host:5432/mydb",
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}

			return runSetEnvCmd(fleetID, args)
		},
	}
	setEnvCmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")

	return setEnvCmd
}

func runSetEnvCmd(fleetID string, args []string) (err error) {
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

	_, err = client.SetEnvironmentVariables(namespace, fleetID, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Environment variable set.")

	return nil
}
