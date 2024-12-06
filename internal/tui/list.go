package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	title       string
	description string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.description }
func (i ListItem) FilterValue() string { return i.title }

type listModel struct {
	list  list.Model
	style lipgloss.Style
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := m.style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m listModel) View() string {
	return m.style.Render(m.list.View())
}

func List(title string, items []list.Item) error {
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = title

	model := listModel{
		list:  list,
		style: lipgloss.NewStyle().Padding(1, 2),
	}

	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}

	return nil
}
