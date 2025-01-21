package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/env"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/cli/tui"
	"github.com/valyentdev/valyent.go"
)

func newDeployCmd() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your project to Valyent",
		RunE: func(cmd *cobra.Command, args []string) error {
			organization, err := cmd.Flags().GetString("organization")
			if err != nil {
				return err
			}
			fleet, err := cmd.Flags().GetString("fleet")
			if err != nil {
				return err
			}
			noBuild, err := cmd.Flags().GetBool("no-build")
			if err != nil {
				return err
			}
			if err := runDeployCmd(organization, fleet, noBuild); err != nil {
				exit.WithError(err)
			}
			return nil
		},
	}
	deployCmd.Flags().String("organization", "", "Organization identifier to use (optional)")
	deployCmd.Flags().String("fleet", "", "Fleet's identifier to use (optional)")
	deployCmd.Flags().Bool("no-build", false, "Do not build the project before deploying (optional)")

	return deployCmd
}

func runDeployCmd(organization, fleet string, noBuild bool) error {
	// Retrieve the project configuration from the `valyent.json`.
	cfg, _ := config.RetrieveProjectConfiguration()
	if fleet == "" {
		fleet = cfg.FleetID
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	var tarball io.ReadCloser
	if !noBuild {
		// Turn current work dir into tarball to upload
		tarball, err = makeTarball()
		if err != nil {
			return fmt.Errorf("failed to prepare tarball of codebase: %v", err)
		}
	}

	if organization == "" {
		namespace, err := config.RetrieveNamespace()
		if err != nil {
			return fmt.Errorf("failed to retrieve namespace: %v", err)
		}
		organization = namespace
	}

	// Create a new deployment by uploading the tarball, ...
	depl, err := client.CreateDeployment(organization, fleet, valyent.CreateDeploymentPayload{
		Machine: cfg.CreateMachinePayload,
	}, tarball)
	if err != nil {
		return fmt.Errorf("failed to create new deployment: %v", err)
	}

	fmt.Println("🎉 Deployment successfully created!")

	baseURL := env.GetVar("VALYENT_API_URL", valyent.DEFAULT_BASE_URL)
	deploymentsPath := fmt.Sprintf(
		"/organizations/%s/applications/%s/deployments",
		organization, fleet,
	)
	deploymentsURL := baseURL + deploymentsPath
	fmt.Printf("You can monitor it at this address: %s\n", deploymentsURL)

	if !noBuild {
		tui.StreamMachineLogs(context.Background(), client, valyent.LogStreamOptions{
			CustomPath: deploymentsPath + "/" + depl.ID + "/builder/logs",
		})
	}

	return nil
}

func makeTarball() (io.ReadCloser, error) {
	paths := listIgnorePaths()

	ar, err := archive.TarWithOptions(".", &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: paths,
	})
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func listIgnorePaths() []string {
	paths := []string{
		"./.git",
		"./node_modules",
		".git/*",
		".git/**/*",
		"node_modules/*",
		"node_modules/**/*",
	}

	readFile, err := os.Open(".dockerignore")
	if err != nil {
		return paths
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		paths = append(paths, fileScanner.Text())
	}

	readFile.Close()

	return paths
}
