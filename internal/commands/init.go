package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/auth"
	"github.com/valyentdev/cli/internal/tui"
	"github.com/valyentdev/cli/pkg/config"
)

func newInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a Valyent configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}
			return runInitCmd(fleetID)
		},
	}

	initCmd.Flags().StringP("fleet", "f", "", "Fleet's identifier to use for initialization (optional)")

	return initCmd
}

func runInitCmd(fleetID string) (err error) {
	// We avoid overwriting an existing configuration file.
	if config.IsAlreadyInitialized() {
		return errors.New("configuration file already exists")
	}

	// If a user already passed a fleet ID through an argument,
	// pass the next steps to initialize the configuration file.
	if fleetID != "" {
		goto init
	}

	// Check that the user is authenticated.
	if !auth.IsLoggedIn() {
		return errors.New("user is not logged in")
	}

	// Select or create fleet.
	fleetID, err = tui.SelectOrCreateFleet()
	if err != nil {
		return
	}

init:
	err = config.InitializeConfigFile(fleetID)
	if err != nil {
		return
	}

	fmt.Println("Valyent configuration file created.")

	return
}
