package gateways

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/internal/tui"
	"github.com/valyentdev/cli/pkg/http"
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
	gatewayName := ""
	err = huh.
		NewInput().
		Title("Type the name of your gateway:").
		Placeholder("bolero-gateway").
		Value(&gatewayName).
		Run()
	if err != nil {
		return
	}

	// We ask for the target port of the gateway.
	rawGatewayTargetPort := ""
	err = huh.
		NewInput().
		Title("Type the target port of your gateway:").
		Placeholder("8080").
		Value(&rawGatewayTargetPort).
		Run()
	if err != nil {
		return
	}

	gatewayTargetPort, err := strconv.Atoi(rawGatewayTargetPort)
	if err != nil {
		return fmt.Errorf("the gateway target port should be a valid integer")
	}

	type createGatewayOptions struct {
		Name       string `json:"name"`
		Fleet      string `json:"fleet"`
		TargetPort int    `json:"target_port"`
	}

	var respTarget any
	err = http.PerformRequest(
		"POST",
		"/v1/gateways",
		createGatewayOptions{
			Fleet:      fleetID,
			Name:       gatewayName,
			TargetPort: gatewayTargetPort,
		},
		&respTarget,
	)
	if err != nil {
		return fmt.Errorf("failed to create gateway: %v", err)
	}

	fmt.Println("🎉 Gateway successfully created!")

	return
}
