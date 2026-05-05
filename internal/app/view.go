package app

import (
	"dawker/internal/command"
	"dawker/internal/ui"
	"fmt"
	"strings"
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

	if m.Loading {
		list = ui.ItemStyle.Render("⏳ Loading resources...")
	} else if m.Err != nil {
		list = ui.StatusExitedStyle.Render("✖ " + m.Err.Error())
	} else {
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
	}

	content := lipgloss.JoinVertical(lipgloss.Left, tabs, "", list)

	w := (m.Width / 2) - 4
	if w < 30 {
		w = 30
	}

	return ui.PanelStyle.
		Width(w).
		Height(m.Height - 6).
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
	if len(m.Containers) == 0 {
		return ui.EmptyStateStyle.Render("No containers found. Press r to refresh.")
	}

	var sb strings.Builder
	for i, c := range m.Containers {
		statusStyle := ui.StatusRunningStyle
		statusText := strings.ToLower(c.Status)
		if strings.Contains(statusText, "exited") || strings.Contains(statusText, "dead") {
			statusStyle = ui.StatusExitedStyle
		}

		status := statusStyle.Render(fmt.Sprintf("(%s)", c.Status))
		line := fmt.Sprintf("%s  %s %s", c.ID, c.Name, status)

		if i == m.SelectedIndexes[TabContainers] {
			sb.WriteString(ui.SelectedItemStyle.Render("▸ "+line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render("  "+line) + "\n")
		}
	}
	return sb.String()
}

func (m Model) renderImages() string {
	if len(m.Images) == 0 {
		return ui.EmptyStateStyle.Render("No images found. Press r to refresh.")
	}

	var sb strings.Builder
	for i, img := range m.Images {
		line := fmt.Sprintf("%s:%s  [%s]", img.Repository, img.Tag, img.Size)

		if i == m.SelectedIndexes[TabImages] {
			sb.WriteString(ui.SelectedItemStyle.Render("▸ "+line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render("  "+line) + "\n")
		}
	}
	return sb.String()
}

func (m Model) renderVolumes() string {
	if len(m.Volumes) == 0 {
		return ui.EmptyStateStyle.Render("No volumes found. Press r to refresh.")
	}

	var sb strings.Builder
	for i, v := range m.Volumes {
		line := fmt.Sprintf("%s  (%s)", v.Name, v.Mountpoint)

		if i == m.SelectedIndexes[TabVolumes] {
			sb.WriteString(ui.SelectedItemStyle.Render("▸ "+line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render("  "+line) + "\n")
		}
	}
	return sb.String()
}

func (m Model) renderNetworks() string {
	if len(m.Networks) == 0 {
		return ui.EmptyStateStyle.Render("No networks found. Press r to refresh.")
	}

	var sb strings.Builder
	for i, n := range m.Networks {
		line := fmt.Sprintf("%s  (%s)", n.Name, n.Driver)

		if i == m.SelectedIndexes[TabNetworks] {
			sb.WriteString(ui.SelectedItemStyle.Render("▸ "+line) + "\n")
		} else {
			sb.WriteString(ui.ItemStyle.Render("  "+line) + "\n")
		}
	}
	return sb.String()
}

func (m Model) renderRightPanel() string {
	w := (m.Width / 2) - 4
	if w < 30 {
		w = 30
	}

	panelHeight := m.Height - 6
	visibleLines := panelHeight - 4
	if visibleLines < 5 {
		visibleLines = 5
	}

	var content string
	titleStr := "DETAILS / LOGS"

	if m.IsStats {
		titleStr = "REAL-TIME STATS"
		content = fmt.Sprintf(
			"%s\n\n%s\n%s\n%s\n%s",
			ui.TitleStyle.Render("Container: "+m.CurrentID),
			fmt.Sprintf("  CPU:    %.2f%%", m.CurrentStats.CPUPercentage),
			fmt.Sprintf("  Memory: %.1fMB / %.1fMB (%.2f%%)", m.CurrentStats.MemoryUsage/1024/1024, m.CurrentStats.MemoryLimit/1024/1024, m.CurrentStats.MemoryPercentage),
			fmt.Sprintf("  Net IO: %s", m.CurrentStats.NetIO),
			fmt.Sprintf("  Blk IO: %s", m.CurrentStats.BlockIO),
		)
	} else if m.Streaming {
		if len(m.Logs) == 0 {
			content = ui.EmptyStateStyle.Render("⏳ Waiting for logs...")
		} else {
			start := 0
			if len(m.Logs) > visibleLines {
				start = len(m.Logs) - visibleLines
			}

			var sb strings.Builder
			for _, line := range m.Logs[start:] {
				if len(line) > w-4 {
					line = line[:w-7] + "..."
				}
				sb.WriteString(line + "\n")
			}
			content = sb.String()
		}
	} else if m.Err != nil {
		content = ui.StatusExitedStyle.Render("✖ Docker error: " + m.Err.Error())
	} else if !m.ShowDetails {
		content = ui.EmptyStateStyle.Render("Press enter on a container to stream logs.")
	} else {
		idx := m.SelectedIndexes[m.ActiveTab]
		switch m.ActiveTab {
		case TabContainers:
			if idx < len(m.Containers) {
				c := m.Containers[idx]
				content = fmt.Sprintf(
					"%s\n%s\n%s",
					ui.TitleStyle.Render("Container: "+c.Name),
					fmt.Sprintf("  ID:     %s", c.ID),
					fmt.Sprintf("  Status: %s", c.Status),
				)
			}
		case TabImages:
			if idx < len(m.Images) {
				img := m.Images[idx]
				content = fmt.Sprintf(
					"%s\n%s\n%s",
					ui.TitleStyle.Render("Image: "+img.Repository),
					fmt.Sprintf("  Tag:  %s", img.Tag),
					fmt.Sprintf("  Size: %s", img.Size),
				)
			}
		case TabVolumes:
			if idx < len(m.Volumes) {
				v := m.Volumes[idx]
				content = fmt.Sprintf(
					"%s\n%s",
					ui.TitleStyle.Render("Volume: "+v.Name),
					fmt.Sprintf("  Mount: %s", v.Mountpoint),
				)
			}
		case TabNetworks:
			if idx < len(m.Networks) {
				n := m.Networks[idx]
				content = fmt.Sprintf(
					"%s\n%s",
					ui.TitleStyle.Render("Network: "+n.Name),
					fmt.Sprintf("  Driver: %s", n.Driver),
				)
			}
		}
	}

	title := ui.TitleStyle.Render(titleStr)
	body := lipgloss.JoinVertical(lipgloss.Left, title, "", content)

	return ui.PanelStyle.
		Width(w).
		Height(panelHeight).
		Render(body)
}

func (m Model) renderFooter() string {
	if m.CommandMode {
		indicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true).
			Render("-- COMMAND --")

		cursor := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Render("█")

		prompt := ":" + m.CommandInput + cursor

		suggestions := ""
		if len(m.Suggestions) > 0 {
			suggestions = "  " + lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render(strings.Join(m.Suggestions, " | "))
		}

		footer := lipgloss.JoinHorizontal(lipgloss.Bottom, indicator, " ", prompt, suggestions)

		if m.CommandError != "" {
			footer += "  " + ui.StatusExitedStyle.Render("✖ "+m.CommandError)
		}

		preview := m.renderCommandPreview()
		if preview != "" {
			return ui.FooterStyle.Width(m.Width).Render(preview + "\n" + footer)
		}

		return ui.FooterStyle.Width(m.Width).Render(footer)
	}

	helpText := "j/k: navigate │ enter: select │ 1-4: tabs │ r: refresh │ : command │ q: quit"
	if m.CommandError != "" {
		helpText = ui.StatusExitedStyle.Render("✖ " + m.CommandError)
	} else if m.CommandResult != "" {
		helpText = ui.StatusRunningStyle.Render("✔ " + m.CommandResult)
	} else if m.Loading {
		helpText = "⏳ Loading..."
	}
	return ui.FooterStyle.Width(m.Width).Render(helpText)
}

func (m Model) renderCommandPreview() string {
	if m.CommandInput == "" {
		return ""
	}

	cmd := command.Parse(m.CommandInput)
	if cmd.Name == "" {
		return ""
	}

	if !command.IsValidCommand(cmd.Name) {
		if suggestion := command.ClosestCommand(cmd.Name); suggestion != "" {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("208")).
				Italic(true).
				Render(fmt.Sprintf("Did you mean: %s?", suggestion))
		}
		return ""
	}

	var names []string

	if len(cmd.Args) == 0 {
		idx := m.SelectedIndexes[m.ActiveTab]
		switch m.ActiveTab {
		case TabContainers:
			if idx < len(m.Containers) {
				names = []string{m.Containers[idx].Name}
			}
		case TabImages:
			if idx < len(m.Images) {
				names = []string{m.Images[idx].Repository}
			}
		case TabVolumes:
			if idx < len(m.Volumes) {
				names = []string{m.Volumes[idx].Name}
			}
		case TabNetworks:
			if idx < len(m.Networks) {
				names = []string{m.Networks[idx].Name}
			}
		}
	} else {
		switch cmd.Name {
		case "rmi":
			for _, i := range command.MatchImages(cmd.Args[0], m.Images) {
				names = append(names, i.Repository)
			}
		case "rmv":
			for _, v := range command.MatchVolumes(cmd.Args[0], m.Volumes) {
				names = append(names, v.Name)
			}
		case "rmn":
			for _, n := range command.MatchNetworks(cmd.Args[0], m.Networks) {
				names = append(names, n.Name)
			}
		case "prune":
			names = []string{"all unused resources"}
		default:
			for _, t := range command.ResolveTargets(cmd.Args, m.Containers) {
				names = append(names, t.Name)
			}
		}
	}

	if len(names) == 0 {
		return ""
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true).
		Render(fmt.Sprintf("Will %s: %s", cmd.Name, strings.Join(names, ", ")))
}
