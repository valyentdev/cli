package commands

import "github.com/spf13/cobra"

func newDeployCmd() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your project to Valyent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployCmd()
		},
	}

	return deployCmd
}

func runDeployCmd() (err error) {
	return
}
