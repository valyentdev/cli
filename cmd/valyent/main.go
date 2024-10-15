package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/valyentdev/cli/internal/commands"
	"github.com/valyentdev/cli/pkg/exit"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	rootCmd := commands.NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		exit.WithError(err)
	}
}
