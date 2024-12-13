package fleets

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/ravel/api"
)

func newCreateFleetCmd() *cobra.Command {
	createFleetCmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			return runCreateFleetCmd(name)
		},
	}

	createFleetCmd.Flags().StringP("name", "f", "", "The name of your fleet (optional)")

	return createFleetCmd
}

func runCreateFleetCmd(name string) error {
	// We ask for the name of the fleet, if not already specified through a flag.
	if name == "" {
		err := huh.
			NewInput().
			Title("Type the name of your fleet:").
			Placeholder("bolero").
			Value(&name).
			Run()
		if err != nil {
			return fmt.Errorf("failed to retrive name for fleet: %v", err)
		}
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking for the creation of a fleet.
	fleet, err := client.CreateFleet(api.CreateFleetPayload{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to create fleet: %v", err)
	}

	fmt.Printf("Fleet successfully created with id \"%s\"!", fleet.Id)

	return nil
}
