package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/pkg/config"
	"github.com/valyentdev/cli/pkg/http"
)

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

func runDeployCmd() error {
	cfg, err := config.RetrieveProjectConfiguration()
	if err != nil {
		return err
	}

	err = http.PerformRequest(
		"POST",
		"/v1/fleets/"+cfg.FleetID+"/machines",
		cfg,
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println("🎉 Machine successfully created!")

	return nil
}
