package main

import (
	"context"
	"os"
	"os/signal"

	commands "github.com/valyentdev/cli/commands"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	rootCmd := commands.NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
