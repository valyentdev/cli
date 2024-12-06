package machines

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/tui"
	"github.com/valyentdev/cli/pkg/http"
	"github.com/valyentdev/ravel/api"
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

	logEntries := []api.LogEntry{}
	err = http.PerformRequest(
		"GET",
		fmt.Sprintf("/v1/fleets/%s/machines/%s/logs", fleetID, machineID),
		nil, &logEntries,
	)
	if err != nil {
		return err
	}

	for _, logEntry := range logEntries {
		if logEntry.Level == "info" {
			logger := log.New(os.Stdout)
			logger.SetReportTimestamp(true)
			logger.SetTimeFunction(func(t time.Time) time.Time {
				return time.Unix(logEntry.Timestamp, 0)
			})
			logger.Info(logEntry.Message)
		}
	}

	return nil
}
