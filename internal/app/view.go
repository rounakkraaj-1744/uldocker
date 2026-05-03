package app

import (
	"fmt"
	"strings"

	"uldocker/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Initializing..."
	}

	leftPanel := m.renderLeftPanel()
	rightPanel := m.renderRightPanel()

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, footer)
}

func (m Model) renderLeftPanel() string {
	tabs := m.renderTabs()
	list := ""

	switch m.ActiveTab {
	case TabContainers:
		list = m.renderContainers()
	case TabImages:
		list = m.renderImages()
	case TabVolumes:
		list = m.renderVolumes()
	case TabNetworks:
		list = m.renderNetworks()
	}

	content := lipgloss.JoinVertical(lipgloss.Left, tabs, "", list)
	
	// Calculate dynamic width based on terminal size, ensure minimum width
	w := (m.Width / 2) - 4
	if w < 30 {
		w = 30
	}
	
	return ui.PanelStyle.
		Width(w).
		Height(m.Height - 6). // account for footer and margins
		Render(content)
}

func (m Model) renderTabs() string {
	tabs := []string{"[1] Containers", "[2] Images", "[3] Volumes", "[4] Networks"}
	var renderedTabs []string

	for i, t := range tabs {
		if Tab(i) == m.ActiveTab {
			renderedTabs = append(renderedTabs, ui.ActiveTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, ui.TabStyle.Render(t))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m Model) renderContainers() string {
	var sb strings.Builder
	for i, c := range m.Containers {
		statusStyle := ui.StatusRunningStyle
		if c.Status == "exited" {
			statusStyle = ui.StatusExitedStyle
		}

		status := statusStyle.Render(fmt.Sprintf("(%s)", c.Status))
		line := fmt.Sprintf("%s - %s %s", c.ID, c.Name, status)

		if i == m.SelectedIndexes[TabContainers] {
			sb.WriteString(ui.SelectedItemStyle.Render(line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render(line) + "\n")
		}
	}
	if len(m.Containers) == 0 {
		sb.WriteString(ui.ItemStyle.Render("No containers found."))
	}
	return sb.String()
}

func (m Model) renderImages() string {
	var sb strings.Builder
	for i, img := range m.Images {
		line := fmt.Sprintf("%s:%s [%s]", img.Repository, img.Tag, img.Size)

		if i == m.SelectedIndexes[TabImages] {
			sb.WriteString(ui.SelectedItemStyle.Render(line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render(line) + "\n")
		}
	}
	if len(m.Images) == 0 {
		sb.WriteString(ui.ItemStyle.Render("No images found."))
	}
	return sb.String()
}

func (m Model) renderVolumes() string {
	var sb strings.Builder
	for i, v := range m.Volumes {
		line := fmt.Sprintf("%s (%s)", v.Name, v.Mountpoint)

		if i == m.SelectedIndexes[TabVolumes] {
			sb.WriteString(ui.SelectedItemStyle.Render(line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render(line) + "\n")
		}
	}
	if len(m.Volumes) == 0 {
		sb.WriteString(ui.ItemStyle.Render("No volumes found."))
	}
	return sb.String()
}

func (m Model) renderNetworks() string {
	var sb strings.Builder
	for i, n := range m.Networks {
		line := fmt.Sprintf("%s (%s driver)", n.Name, n.Driver)

		if i == m.SelectedIndexes[TabNetworks] {
			sb.WriteString(ui.SelectedItemStyle.Render(line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render(line) + "\n")
		}
	}
	if len(m.Networks) == 0 {
		sb.WriteString(ui.ItemStyle.Render("No networks found."))
	}
	return sb.String()
}

func (m Model) renderRightPanel() string {
	w := (m.Width / 2) - 4
	if w < 30 {
		w = 30
	}

	content := ""

	if !m.ShowDetails {
		content = "Select an item to view details"
	} else {
		idx := m.SelectedIndexes[m.ActiveTab]
		switch m.ActiveTab {
		case TabContainers:
			if len(m.Containers) > idx {
				content = fmt.Sprintf("Viewing logs for %s", m.Containers[idx].Name)
			}
		case TabImages:
			if len(m.Images) > idx {
				content = fmt.Sprintf("Image details: %s:%s", m.Images[idx].Repository, m.Images[idx].Tag)
			}
		case TabVolumes:
			if len(m.Volumes) > idx {
				content = fmt.Sprintf("Volume details: %s", m.Volumes[idx].Name)
			}
		case TabNetworks:
			if len(m.Networks) > idx {
				content = fmt.Sprintf("Network details: %s", m.Networks[idx].Name)
			}
		}
	}

	title := ui.TitleStyle.Render("DETAILS / LOGS")
	body := lipgloss.JoinVertical(lipgloss.Left, title, "", content)

	return ui.PanelStyle.
		Width(w).
		Height(m.Height - 6).
		Render(body)
}

func (m Model) renderFooter() string {
	helpText := "j/k: move | enter: select | 1-4: switch tabs | : command | q: quit"
	return ui.FooterStyle.Width(m.Width).Render(helpText)
}