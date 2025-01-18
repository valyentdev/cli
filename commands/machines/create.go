package machines

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
	"github.com/valyentdev/ravel/api"
)

func newCreateMachineCmd() *cobra.Command {
	createMachineCmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonFilePath, err := cmd.Flags().GetString("json-file")
			if err != nil {
				return err
			}

			if err := runCreateMachineCmd(jsonFilePath); err != nil {
				exit.WithError(err)
			}

			return nil
		},
	}

	createMachineCmd.Flags().
		StringP("json-file", "j", "", "The path to your JSON machine configuration file")

	return createMachineCmd
}

func runCreateMachineCmd(jsonFilePath string) error {
	// Check that our JSON config file is correctly set.
	if jsonFilePath == "" {
		return errors.New("the json config file path should not be empty")
	}

	// Load the JSON machine config file.
	cfg := api.CreateMachinePayload{}
	contents, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read json config file: %v", err)
	}
	if err := json.Unmarshal(contents, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal json contents into struct: %v", err)
	}

	// Select fleet.
	fleetID, err := tui.SelectFleet()
	if err != nil {
		return fmt.Errorf("failed to select fleet: %v", err)
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking for machine creation.
	machine, err := client.CreateMachine(fleetID, cfg)
	if err != nil {
		return fmt.Errorf("failed to create machine: %v", err)
	}

	fmt.Printf("🚀 Machine successfully created with id: \"%s\".\n", machine.Id)

	return nil
}
