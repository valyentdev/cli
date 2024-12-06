package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FancySelectItem struct {
	title       string
	description string
}

func (i FancySelectItem) Title() string       { return i.title }
func (i FancySelectItem) Description() string { return i.description }
func (i FancySelectItem) FilterValue() string { return i.title }

type fancySelectModel struct {
	list     list.Model
	style    lipgloss.Style
	quitting bool
	selected string
}

func (m *fancySelectModel) Init() tea.Cmd {
	return nil
}

func (m *fancySelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c": // Quit the program
			m.quitting = true
			return m, tea.Quit
		case "enter": // Select the current item
			if selectedItem, ok := m.list.SelectedItem().(FancySelectItem); ok {
				m.selected = selectedItem.title
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := m.style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *fancySelectModel) View() string {
	if m.quitting {
		return "Goodbye!"
	}
	if m.selected != "" {
		return fmt.Sprintf("You selected: %s", m.selected)
	}
	return m.style.Render(m.list.View())
}

// FancySelect displays a list of items and allows the user to select one.
func FancySelect(title string, items []list.Item) (string, error) {
	// Create the list model
	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = title

	// Wrap it in our fancySelectModel
	model := &fancySelectModel{
		list:  listModel,
		style: lipgloss.NewStyle().Padding(1, 2),
	}

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return "", err
	}

	return model.selected, nil
}
