package machines

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/logs"
	"github.com/valyentdev/cli/tui"
	"github.com/valyentdev/valyent.go"
)

func newLogsCmd() *cobra.Command {
	var follow bool

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Display logs from a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogsCmd(follow)
		},
	}

	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow the logs")
	return logsCmd
}

func runLogsCmd(follow bool) error {
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

	if follow {
		return logs.Stream(client, valyent.LogStreamOptions{
			FleetID:   fleetID,
			MachineID: machineID,
		})
	}

	return displayLogs(client, fleetID, machineID)
}

func displayLogs(client *valyent.Client, fleetID, machineID string) error {
	logEntries, err := client.GetLogs(fleetID, machineID)
	if err != nil {
		return fmt.Errorf("failed to fetch logs from the API: %v", err)
	}

	for _, logEntry := range logEntries {
		fmt.Println(logs.FormatEntry(logEntry))
	}

	return nil
}
