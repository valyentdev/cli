package commands

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/auth"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/env"
	"github.com/valyentdev/cli/tui"
)

func newOpenCmd() *cobra.Command {
	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Open application in the browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOpenCmd()
		},
	}

	return openCmd
}

func runOpenCmd() (err error) {
	// Check that the user is authenticated.
	if !auth.IsLoggedIn() {
		return errors.New("user is not logged in")
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	fleetID, err := tui.SelectFleet()
	if err != nil {
		return err
	}

	fleet, err := client.GetFleet(fleetID)
	if err != nil {
		return err
	}

	gtw, err := tui.SelectGateway(fleet)
	if err != nil {
		return fmt.Errorf("failed to select gateway: %v", err)
	}

	namespace, err := config.RetrieveNamespace()
	if err != nil {
		return fmt.Errorf("failed to retrieve namespace: %v", namespace)
	}

	host := env.GetVar("VALYENT_WILDCARD_DOMAIN", "valyent.app")
	url := "https://" + gtw.Name + "-" + namespace + "." + host
	if err := openBrowser(url); err != nil {
		fmt.Println(url)
	}

	return
}

func openBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return
}
