package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	api "github.com/valyentdev/valyent.go"
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
	// Retrieve the project configuration from the `valyent.json`.
	cfg, err := config.RetrieveProjectConfiguration()
	if err != nil {
		return err
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Turn current work dir into tarball to upload
	tarball, err := makeTarball()
	if err != nil {
		return fmt.Errorf("failed to prepare tarball of codebase: %v", err)
	}

	namespace, err := config.RetrieveNamespace()
	if err != nil {
		return fmt.Errorf("failed to retrieve namespace: %v", err)
	}

	// Create a new deployment by uploading the tarball, ...
	_, err = client.CreateDeployment(namespace, cfg.FleetID, api.CreateDeploymentPayload{
		Machine: cfg.CreateMachinePayload,
	}, tarball)
	if err != nil {
		return fmt.Errorf("failed to create new deployment: %v", err)
	}

	fmt.Println("🎉 Deployment successfully created!")

	deploymentURL := "https://console.valyent.cloud/organizations/" +
		namespace + "/applications/+" + cfg.FleetID + "/deployments"
	fmt.Printf("You can monitor it at this address: %s\n", deploymentURL)

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
