package fleets

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/tui"
)

func newDeleteFleetCmd() *cobra.Command {
	deleteFleetCmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"remove"},
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
	// If the fleet if not specified by a flag, we let the user visually select it.
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

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking for fleet deletion.
	if err = client.DeleteFleet(fleetID); err != nil {
		return err
	}

	fmt.Println("Fleet successfully deleted.")

	return
}
