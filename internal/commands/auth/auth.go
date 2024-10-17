package auth

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use: "auth",
	}

	authCmd.AddCommand(newLoginCmd())
	authCmd.AddCommand(newLogoutCmd())

	return authCmd
}
