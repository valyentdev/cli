package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/pkg/config"
)

func newLogoutCmd() *cobra.Command {
	logoutCmd := &cobra.Command{
		Use: "logout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogoutCmd()
		},
	}

	return logoutCmd
}

func runLogoutCmd() (err error) {
	err = config.RemoveConfigFile()
	if err != nil {
		return
	}

	fmt.Println("Successfully logged out.")

	return
}
