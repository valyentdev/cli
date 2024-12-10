package fleets

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/api"
	"github.com/valyentdev/cli/tui"
)

func newDeleteFleetCmd() *cobra.Command {
	deleteFleetCmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}
			return runDeleteFleetCmd(fleetID)
		},
	}

	deleteFleetCmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")

	return deleteFleetCmd
}

func runDeleteFleetCmd(fleetID string) (err error) {
	if fleetID == "" {
		fleetID, err = tui.SelectFleet()
		if err != nil {
			return err
		}
	}

	confirmed := false
	err = huh.NewConfirm().
		Title("Are you sure you want to delete this fleet?").
		Value(&confirmed).
		Run()
	if err != nil {
		return err
	}
	if !confirmed {
		return
	}

	if err = api.DeleteFleet(fleetID); err != nil {
		return err
	}

	fmt.Println("Fleet successfully deleted.")

	return
}
