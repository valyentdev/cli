package tui

import (
	"fmt"
	stdHTTP "net/http"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/valyentdev/cli/internal/api"
	"github.com/valyentdev/cli/pkg/http"
	ravelAPI "github.com/valyentdev/ravel/api"
)

// SelectFleet prompts the user to select an existing fleet from the list.
func SelectFleet() (fleetID string, err error) {
	// Retrieve fleets from the API.
	fleets, err := api.GetFleets()
	if err != nil {
		return
	}

	return SelectFleetWithFleets(fleets)
}

// SelectFleetWithFleets prompts the user to select an existing fleet from the list.
func SelectFleetWithFleets(fleets []ravelAPI.Fleet) (fleetID string, err error) {
	// Compute select options.
	opts := []huh.Option[string]{}
	for _, fleet := range fleets {
		opts = append(opts, huh.NewOption(fleet.Name, fleet.Id))
	}

	if len(fleets) > 0 {
		fleetID = fleets[0].Id
	}

	// Ask the user to select a fleet.
	err = huh.
		NewSelect[string]().
		Title("Pick a fleet").
		Options(opts...).
		Value(&fleetID).
		Height(10).
		Run()

	return
}

// SelectOrCreateFleet lets the user select or create a fleet.
func SelectOrCreateFleet() (fleetID string, err error) {
	// Retrieve fleets from the API.
	fleets, err := api.GetFleets()
	if err != nil {
		return
	}

	// Append an empty fleet to serve as a creation option.
	fleets = append(fleets, ravelAPI.Fleet{
		Id:   "",
		Name: "[+] Create a new fleet",
	})

	// Attempt to select an existing fleet.
	fleetID, err = SelectFleetWithFleets(fleets)
	if err != nil {
		return
	}

	// Handle fleet creation if no existing fleet is selected.
	if fleetID == "" {
		fleetName := ""
		err = huh.
			NewInput().
			Title("Type the name of your fleet:").
			Placeholder("bolero").
			Value(&fleetName).
			Run()
		if err != nil {
			return
		}

		// Create the new fleet and assign its ID.
		err = spinner.
			New().
			Title("Creating fleet...").
			Action(func() {
				payload := struct {
					Name string `json:"name"`
				}{
					Name: fleetName,
				}
				createFleetResponse := &ravelAPI.Fleet{}
				err = http.PerformRequest(stdHTTP.MethodPost, "/v1/fleets", payload, createFleetResponse)
				if err != nil {
					return
				}
				fleetID = createFleetResponse.Id
			}).
			Run()
		if err != nil {
			return
		}
	}

	return
}

// ListFleets lists the fleets related to the currently authenticated's namespace
// without prompting the user to choose one for further actions.
func ListFleets() error {
	// Retrieve fleets from the API.
	fleets, err := api.GetFleets()
	if err != nil {
		return err
	}

	// Make initial list of items
	items := make([]list.Item, len(fleets))
	for idx, fleet := range fleets {
		items[idx] = ListItem{
			title:       fleet.Name,
			description: fmt.Sprintf("Created at: %s", fleet.CreatedAt.Format("Monday, Jan 2, 2006 at 3:04 PM")),
		}
	}

	return List("List of Fleets", items)
}