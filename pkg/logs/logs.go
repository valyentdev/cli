package logs

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/valyentdev/ravel/api"
	"github.com/valyentdev/valyent.go"
)

func Stream(client *valyent.Client, opts valyent.LogStreamOptions) error {
	// Create a context that we can cancel on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up interrupt handling
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// Start streaming logs
	stream, err := client.StreamLogs(ctx, opts)
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

		fmt.Println(FormatEntry(entry))
	}
}

func FormatEntry(entry api.LogEntry) string {
	timestamp := time.Unix(entry.Timestamp, 0).Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] [%s] %s", timestamp, entry.Level, entry.Message)
}
