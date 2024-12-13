package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
)

func newDeployCmd() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your project to Valyent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployCmd()
		},
	}

	return deployCmd
}

func runDeployCmd() error {
	// Retrieve the project configuration from the `valyent.json`.
	cfg, err := config.RetrieveProjectConfiguration()
	if err != nil {
		return err
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking for machine creation.
	if _, err := client.CreateMachine(cfg.FleetID, cfg.CreateMachinePayload); err != nil {
		return fmt.Errorf("failed to create machine: %v", err)
	}

	fmt.Println("🎉 Machine successfully created!")

	return nil
}
