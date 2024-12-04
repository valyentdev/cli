package commands

import (
	"errors"
	"fmt"
	stdHTTP "net/http"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/auth"
	"github.com/valyentdev/cli/pkg/config"
	"github.com/valyentdev/cli/pkg/http"
	"github.com/valyentdev/ravel/api"
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
	if config.IsAlreadyInitialized() {
		return errors.New("configuration file already exists")
	}

	if fleetID != "" {
		goto init
	}

	if !auth.IsLoggedIn() {
		return errors.New("user is not logged in")
	}

	{
		fleets := &[]api.Fleet{}

		err = http.PerformRequest(stdHTTP.MethodGet, "/v1/fleets", nil, fleets)
		if err != nil {
			return
		}

		opts := []huh.Option[string]{}
		for _, fleet := range *fleets {
			opts = append(opts, huh.NewOption(fleet.Name, fleet.Id))
		}
		opts = append(opts, huh.NewOption("[+] Create a new fleet", ""))

		err = huh.
			NewSelect[string]().
			Options(opts...).
			Value(&fleetID).
			Height(10).
			Run()
		if err != nil {
			return
		}

		if fleetID == "" {
			fleetName := ""

			err = huh.
				NewInput().
				Title("Type the name of your fleet:").
				Placeholder("bolero").
				Value(&fleetName).
				Run()
			if err != nil {
				return err
			}

			err = spinner.
				New().
				Title("Creating fleet...").
				Action(func() {
					payload := struct {
						Name string `json:"name"`
					}{
						Name: fleetName,
					}
					createFleetResponse := &api.Fleet{}
					err = http.PerformRequest(stdHTTP.MethodPost, "/v1/fleets", payload, createFleetResponse)
					if err != nil {
						return
					}
					fleetID = createFleetResponse.Id
				}).
				Run()
			if err != nil {
				return err
			}
		}
	}

init:
	err = config.InitializeConfigFile(fleetID)
	if err != nil {
		return
	}

	fmt.Println("Valyent configuration file created.")

	return
}
