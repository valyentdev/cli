package fleets

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/pkg/http"
)

func newCreateFleetCmd() *cobra.Command {
	createFleetCmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateFleetCmd()
		},
	}

	return createFleetCmd
}

func runCreateFleetCmd() error {
	// We ask for the name of the fleet.
	fleetName := ""
	err := huh.
		NewInput().
		Title("Type the name of your fleet:").
		Placeholder("bolero").
		Value(&fleetName).
		Run()
	if err != nil {
		return fmt.Errorf("failed to retrive name for fleet: %v", err)
	}

	type createFleetOptions struct {
		Name string `json:"name"`
	}

	err = http.PerformRequest(
		"POST",
		"/v1/fleets",
		createFleetOptions{
			Name: fleetName,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create fleet: %v", err)
	}

	fmt.Println("🎉 Fleet successfully created!")

	return nil
}
