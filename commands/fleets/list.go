package fleets

import (
	"github.com/spf13/cobra"
	tui "github.com/valyentdev/cli/tui"
)

func newListFleetsCmd() *cobra.Command {
	listFleetsCmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListFleetsCmd()
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
