package machines

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/tui"
	ravelAPI "github.com/valyentdev/ravel/api"
	api "github.com/valyentdev/valyent.go"
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
		return streamLogs(client, fleetID, machineID)
	}

	return displayLogs(client, fleetID, machineID)
}

func displayLogs(client *api.Client, fleetID, machineID string) error {
	logEntries, err := client.GetLogs(fleetID, machineID)
	if err != nil {
		return fmt.Errorf("failed to fetch logs from the API: %v", err)
	}

	for _, logEntry := range logEntries {
		fmt.Println(formatLogEntry(logEntry))
	}

	return nil
}

func streamLogs(client *api.Client, fleetID, machineID string) error {
	// Create a context that we can cancel on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up interrupt handling
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// Start streaming logs
	stream, err := client.StreamLogs(ctx, api.LogStreamOptions{
		FleetID:   fleetID,
		MachineID: machineID,
		Follow:    true,
	})
	if err != nil {
		return fmt.Errorf("failed to start log stream: %v", err)
	}
	defer stream.Close()

	// Read and print logs until interrupted or error
	for {
		entry, ok := stream.Next()
		if !ok {
			if err := stream.Err(); err != nil {
				return fmt.Errorf("stream error: %v", err)
			}
			return nil // Clean exit (context cancelled)
		}

		fmt.Println(formatLogEntry(entry))
	}
}

func formatLogEntry(entry ravelAPI.LogEntry) string {
	timestamp := time.Unix(entry.Timestamp, 0).Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] [%s] %s", timestamp, entry.Level, entry.Message)
}
