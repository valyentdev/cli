package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/valyentdev/cli/http"
	ravelAPI "github.com/valyentdev/ravel/api"
)

func ListMachines(fleets []ravelAPI.Fleet, fleetID string) error {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Retrieve machines from the API.
	machines, err := client.GetMachines(fleetID)
	if err != nil {
		return err
	}

	// Make initial list of items
	items := make([]list.Item, len(machines))
	for idx, machine := range machines {
		items[idx] = ListItem{
			title:       machine.Id,
			description: getDescriptionForMachine(machine),
		}
	}

	return List("List of Machines", items)
}

func getDescriptionForMachine(machine ravelAPI.Machine) string {
	// Compute image string.
	var image string = machine.Config.Image

	splittedImage := strings.Split(machine.Config.Image, "@")
	if len(splittedImage) > 1 {
		image = splittedImage[0]
	}

	return fmt.Sprintf(
		"Image: %s | Region: %s | Status: %s",
		image, machine.Region, getStatusString(string(machine.Status)),
	)
}

func getStatusString(status string) string {
	switch status {
	case "created":
		return "🆕 Created"
	case "preparing":
		return "⚙️ Preparing"
	case "starting":
		return "🚀 Starting"
	case "running":
		return "🏃‍ Running"
	case "stopping":
		return "🛑 Stopping"
	case "stopped":
		return "⛔ Stopped"
	case "destroying":
		return "💣 Destroying"
	case "destroyed":
		return "💥 Destroyed"
	default:
		return ""
	}
}

func SelectMachine(fleetID string) (string, error) {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Retrieve machines from the API.
	machines, err := client.GetMachines(fleetID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve machines from the API: %v", err)
	}

	items := []list.Item{}

	for _, machine := range machines {
		items = append(items, FancySelectItem{
			title:       machine.Id,
			description: getDescriptionForMachine(machine),
			value:       machine.Id,
		})
	}

	machineID, err := FancySelect("Select a Machine", items)
	if err != nil {
		return "", err
	}

	if machineID == "" {
		return "", errors.New("no machine has been selected")
	}

	return machineID, nil
}

func ListMachineEvents(fleetID, machineID string) error {
	// Initialize new Valyent API HTTP client.
	client, err := http.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
	}

	// Retrieve events from the API.
	events, err := client.GetMachineEvents(fleetID, machineID)
	if err != nil {
		return err
	}

	// Make initial list of items
	items := make([]list.Item, len(events))
	for idx, event := range events {
		items[idx] = ListItem{
			title:       string(event.Status),
			description: getMachineEventDescription(event),
		}

	}

	return List("List of Machine Events", items)
}

func getMachineEventDescription(event ravelAPI.MachineEvent) string {
	timestamp := event.Timestamp.Format("2006/01/02 15:04")
	return fmt.Sprintf("Timestamp: %s", timestamp)
}
