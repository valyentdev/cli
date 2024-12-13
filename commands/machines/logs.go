package machines

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/tui"
)

func newLogsCmd() *cobra.Command {
	logsCmd := &cobra.Command{
		Use: "logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogsCmd()
		},
	}

	return logsCmd
}

func runLogsCmd() error {
	// Allow user to select fleet.
	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	// Allow user to select a machine related to a specific fleet.
	machineID, err := tui.SelectMachine(fleetID)
	if err != nil {
		return err
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Let's fetch logs from the API.
	logEntries, err := client.GetLogs(fleetID, machineID)
	if err != nil {
		return fmt.Errorf("failed to fetch logs from the API: %v", err)
	}

	// For now, we just print each message from log entries.
	for _, logEntry := range logEntries {
		fmt.Println(logEntry.Message)
	}

	return nil
}
