package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "1":
			m.ActiveTab = TabContainers
			m.ShowDetails = false
		case "2":
			m.ActiveTab = TabImages
			m.ShowDetails = false
		case "3":
			m.ActiveTab = TabVolumes
			m.ShowDetails = false
		case "4":
			m.ActiveTab = TabNetworks
			m.ShowDetails = false

		case "j", "down":
			m.moveSelection(1)
			m.ShowDetails = false

		case "k", "up":
			m.moveSelection(-1)
			m.ShowDetails = false

		case "enter":
			m.ShowDetails = true

		case "esc":
			m.ShowDetails = false
		}
	}

	return m, nil
}

func (m *Model) moveSelection(amount int) {
	currentIdx := m.SelectedIndexes[m.ActiveTab]
	newIdx := currentIdx + amount

	var maxIdx int
	switch m.ActiveTab {
	case TabContainers:
		maxIdx = len(m.Containers) - 1
	case TabImages:
		maxIdx = len(m.Images) - 1
	case TabVolumes:
		maxIdx = len(m.Volumes) - 1
	case TabNetworks:
		maxIdx = len(m.Networks) - 1
	}

	if maxIdx < 0 {
		m.SelectedIndexes[m.ActiveTab] = 0
		return
	}

	if newIdx < 0 {
		newIdx = 0
	} else if newIdx > maxIdx {
		newIdx = maxIdx
	}

	m.SelectedIndexes[m.ActiveTab] = newIdx
}