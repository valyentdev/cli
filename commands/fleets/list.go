package fleets

import (
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/pkg/exit"
	tui "github.com/valyentdev/cli/tui"
)

func newListFleetsCmd() *cobra.Command {
	listFleetsCmd := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runListFleetsCmd(); err != nil {
				exit.WithError(err)
			}
		},
	}

	return listFleetsCmd
}

func runListFleetsCmd() error {
	err := tui.ListFleets()
	if err != nil {
		return err
	}

	return nil
}
