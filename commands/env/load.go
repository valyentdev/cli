package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
)

func newLoadEnvCmd() *cobra.Command {
	envLoadCmd := &cobra.Command{
		Use:     "load",
		Short:   "Load environment variables from file",
		Example: "valyent env load .env",
		RunE: func(cmd *cobra.Command, args []string) error {
			fleetID, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}

			return runLoadEnvCmd(fleetID, args)
		},
	}
	envLoadCmd.Flags().StringP("fleet", "f", "", "Fleet's identifier (optional)")

	return envLoadCmd
}

func runLoadEnvCmd(fleetID string, args []string) (err error) {
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

	if len(args) != 1 {
		fmt.Println("Please provide the path to the file from which you want to load the environment variables")
		os.Exit(1)
	}

	vars, err := retrieveEnvironmentVariablesFromFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	_, err = client.SetEnvironmentVariables(namespace, fleetID, vars)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Environment variables loaded.")

	return nil
}

// retrieveEnvironmentVariablesFromFile reads a .env file and returns an array of strings in the format "key=value".
func retrieveEnvironmentVariablesFromFile(path string) ([]string, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var envVariables []string

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 || !strings.Contains(line, "=") {
			continue
		}

		// Add the line to the array
		envVariables = append(envVariables, line)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVariables, nil
}
