package gateways

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/http"
	tui "github.com/valyentdev/cli/tui"
	"github.com/valyentdev/ravel/api"
)

func newCreateGatewayCmd() *cobra.Command {
	createGatewayCmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateGatewayCmd()
		},
	}

	return createGatewayCmd
}

func runCreateGatewayCmd() (err error) {
	// We select the fleet for which we want to create a gateway.
	fleetID, err := tui.SelectFleet()
	if err != nil {
		return
	}

	// We ask for the name of the gateway.
	name := ""
	err = huh.
		NewInput().
		Title("Type the name of your gateway:").
		Placeholder("bolero-gateway").
		Value(&name).
		Run()
	if err != nil {
		return
	}

	// We ask for the target port of the gateway.
	rawTargetPort := ""
	err = huh.
		NewInput().
		Title("Type the target port of your gateway:").
		Placeholder("8080").
		Value(&rawTargetPort).
		Run()
	if err != nil {
		return
	}

	targetPort, err := strconv.Atoi(rawTargetPort)
	if err != nil {
		return fmt.Errorf("the gateway target port should be a valid integer")
	}

	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Call the API asking for gateway creation.
	gtw, err := client.CreateGateway(api.CreateGatewayPayload{
		Fleet:      fleetID,
		Name:       name,
		TargetPort: targetPort,
	})
	if err != nil {
		return fmt.Errorf("failed to create gateway: %v", err)
	}

	fmt.Printf("✅ Gateway successfully created with ID \"%s\".\n", gtw.Id)

	return
}
