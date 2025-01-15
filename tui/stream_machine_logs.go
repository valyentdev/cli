package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/valyentdev/valyent.go"
)

// StreamMachineLogs streams logs using the StreamModel and the LogStream API.
func StreamMachineLogs(ctx context.Context, client *valyent.Client, opts valyent.LogStreamOptions) {
	// Create a new StreamModel for the TUI
	streamModel := NewStreamModel("Streaming Logs...")

	// Variable to handle errors
	var streamErr error

	// Run the TUI model
	streamModel.Run(func(p *tea.Program) {
		// Start the log stream
		logStream, err := client.StreamLogs(ctx, opts)
		if err != nil {
			fmt.Println("Error starting log stream:", err)
			os.Exit(1)
		}
		defer logStream.Close()

		// Read log entries
		for {
			entry, ok := logStream.Next()
			if !ok {
				if err := logStream.Err(); err != nil {
					streamErr = fmt.Errorf("error streaming logs: %w", err)
				}
				p.Quit()
				break
			}

			// Format the log entry for display
			timestamp := time.Unix(entry.Timestamp, 0)
			formattedTime := timestamp.Format("15:04:05") // Display time only
			message := strings.TrimSpace(entry.Message)

			// Send the log entry to the StreamModel
			p.Send(StreamModelResultMsg{
				Title:    message,
				Subtitle: formattedTime,
			})
		}
	})

	// Handle errors after the TUI exits
	if streamErr != nil {
		fmt.Println("🔴 Error streaming logs:", streamErr)
		os.Exit(1)
	}

	fmt.Println("✅ Log streaming completed.")
}
