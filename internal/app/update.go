package app

import (
	"context"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/pkg/stdcopy"
	"uldocker/internal/command"
	"uldocker/internal/docker"
	"uldocker/pkg/types"
)

type resourcesLoadedMsg struct {
	containers []types.Container
	images     []types.Image
	volumes    []types.Volume
	networks   []types.Network
	err        error
}

func loadResources() tea.Msg {
	containers, err := docker.ListContainers()
	if err != nil {
		return resourcesLoadedMsg{err: err}
	}
	images, err := docker.ListImages()
	if err != nil {
		return resourcesLoadedMsg{err: err}
	}
	volumes, err := docker.ListVolumes()
	if err != nil {
		return resourcesLoadedMsg{err: err}
	}
	networks, err := docker.ListNetworks()
	if err != nil {
		return resourcesLoadedMsg{err: err}
	}

	return resourcesLoadedMsg{
		containers: containers,
		images:     images,
		volumes:    volumes,
		networks:   networks,
	}
}

type logStreamStartedMsg struct {
	reader io.Reader
}

type logChunkMsg struct {
	reader io.Reader
	chunk  string
}

type streamEndedMsg struct{}

func startStreamCmd(ctx context.Context, containerID string) tea.Cmd {
	return func() tea.Msg {
		reader, err := docker.StreamLogs(ctx, containerID)
		if err != nil {
			return logChunkMsg{chunk: "ERROR: " + err.Error()}
		}

		pr, pw := io.Pipe()

		go func() {
			defer reader.Close()
			defer pw.Close()

			go func() {
				<-ctx.Done()
				reader.Close()
				pw.Close()
			}()

			_, err := stdcopy.StdCopy(pw, pw, reader)
			if err != nil {
				io.Copy(pw, reader)
			}
		}()

		return logStreamStartedMsg{reader: pr}
	}
}

func readStreamCmd(reader io.Reader) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, 4096)
		n, err := reader.Read(buf)
		if err != nil {
			return streamEndedMsg{}
		}
		return logChunkMsg{
			reader: reader,
			chunk:  string(buf[:n]),
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case resourcesLoadedMsg:
		m.Loading = false
		if msg.err != nil {
			m.Err = msg.err
			return m, nil
		}
		m.Containers = msg.containers
		m.Images = msg.images
		m.Volumes = msg.volumes
		m.Networks = msg.networks
		return m, nil

	case tea.KeyMsg:
		if m.CommandResult != "" {
			m.CommandResult = ""
		}

		if m.CommandMode {
			switch msg.Type {
			case tea.KeyEnter:
				return m.executeCommand()

			case tea.KeyEsc:
				m.CommandMode = false
				m.Suggestions = nil
				return m, nil

			case tea.KeyUp:
				if len(m.History) > 0 {
					if m.HistoryIndex == -1 {
						m.HistoryIndex = len(m.History) - 1
					} else if m.HistoryIndex > 0 {
						m.HistoryIndex--
					}
					m.CommandInput = m.History[m.HistoryIndex]
					m.updateSuggestions()
				}
				return m, nil

			case tea.KeyDown:
				if len(m.History) > 0 && m.HistoryIndex != -1 {
					if m.HistoryIndex < len(m.History)-1 {
						m.HistoryIndex++
						m.CommandInput = m.History[m.HistoryIndex]
					} else {
						m.HistoryIndex = -1
						m.CommandInput = ""
					}
					m.updateSuggestions()
				}
				return m, nil

			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.CommandInput) > 0 {
					m.CommandInput = m.CommandInput[:len(m.CommandInput)-1]
				}
				m.updateSuggestions()

			default:
				m.CommandInput += msg.String()
				m.updateSuggestions()
			}

			return m, nil
		}

		switch msg.String() {
		case ":":
			m.CommandMode = true
			m.CommandInput = ""
			m.CommandError = ""
			return m, nil

		case "esc", "1", "2", "3", "4", "j", "k", "up", "down":
			if msg.String() != "enter" && m.LogsCancel != nil {
				m.LogsCancel()
				m.LogsCancel = nil
				m.Streaming = false
			}
			// fallthrough to handle the actual key
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			m.Loading = true
			return m, loadResources

		case "j", "down":
			m.moveDown()

		case "k", "up":
			m.moveUp()

		case "1":
			m.ActiveTab = TabContainers

		case "2":
			m.ActiveTab = TabImages

		case "3":
			m.ActiveTab = TabVolumes

		case "4":
			m.ActiveTab = TabNetworks

		case "enter":
			if m.ActiveTab == TabContainers && len(m.Containers) > 0 {
				idx := m.SelectedIndexes[TabContainers]
				container := m.Containers[idx]

				if m.LogsCancel != nil {
					m.LogsCancel()
				}

				m.Logs = []string{}
				m.Streaming = true
				m.CurrentID = container.ID

				ctx, cancel := context.WithCancel(context.Background())
				m.LogsCancel = cancel

				return m, startStreamCmd(ctx, container.ID)
			}
			m.ShowDetails = !m.ShowDetails
		}
		return m, nil

	case logStreamStartedMsg:
		return m, readStreamCmd(msg.reader)

	case logChunkMsg:
		lines := strings.Split(string(msg.chunk), "\n")

		if len(m.Logs) == 0 {
			m.Logs = lines
		} else {
			m.Logs[len(m.Logs)-1] += lines[0]
			if len(lines) > 1 {
				m.Logs = append(m.Logs, lines[1:]...)
			}
		}

		if len(m.Logs) > 100 {
			m.Logs = m.Logs[len(m.Logs)-100:]
		}
		return m, readStreamCmd(msg.reader)

	case streamEndedMsg:
		return m, nil
	}

	return m, nil
}

func (m *Model) moveDown() {
	idx := m.SelectedIndexes[m.ActiveTab]
	max := m.getActiveListLength()

	if idx < max-1 {
		m.SelectedIndexes[m.ActiveTab] = idx + 1
	}
}

func (m *Model) moveUp() {
	idx := m.SelectedIndexes[m.ActiveTab]

	if idx > 0 {
		m.SelectedIndexes[m.ActiveTab] = idx - 1
	}
}

func (m *Model) getActiveListLength() int {
	switch m.ActiveTab {
		case TabContainers:
			return len(m.Containers)
		case TabImages:
			return len(m.Images)
		case TabVolumes:
			return len(m.Volumes)
		case TabNetworks:
			return len(m.Networks)
		default:
			return 0
	}
}

func (m Model) executeCommand() (tea.Model, tea.Cmd) {
	cmd := command.Parse(m.CommandInput)
	if cmd.Name == "" {
		m.CommandMode = false
		return m, nil
	}

	var targets []types.Container
	var images  []types.Image
	var volumes []types.Volume
	var networks []types.Network

	// Context Awareness based on active tab
	if len(cmd.Args) == 0 {
		idx := m.SelectedIndexes[m.ActiveTab]
		switch m.ActiveTab {
		case TabContainers:
			if len(m.Containers) > idx {
				targets = []types.Container{m.Containers[idx]}
			}
		case TabImages:
			if len(m.Images) > idx {
				images = []types.Image{m.Images[idx]}
			}
		case TabVolumes:
			if len(m.Volumes) > idx {
				volumes = []types.Volume{m.Volumes[idx]}
			}
		case TabNetworks:
			if len(m.Networks) > idx {
				networks = []types.Network{m.Networks[idx]}
			}
		}
	} else {
		// Smart Resolution
		switch cmd.Name {
		case "rmi":
			images = command.MatchImages(cmd.Args[0], m.Images)
		case "rmv":
			volumes = command.MatchVolumes(cmd.Args[0], m.Volumes)
		case "rmn":
			networks = command.MatchNetworks(cmd.Args[0], m.Networks)
		case "prune":
			// Prune acts on everything
		default:
			targets = command.ResolveTargets(cmd.Args, m.Containers)
		}
	}

	resultMsg, err := command.Execute(cmd, targets, images, volumes, networks)

	if err != nil {
		m.CommandError = err.Error()
		m.CommandResult = ""
	} else {
		m.CommandError = ""
		m.CommandResult = resultMsg
		
		// History
		m.History = append(m.History, m.CommandInput)
		m.HistoryIndex = -1
	}

	m.CommandMode = false
	m.Suggestions = nil
	if err == nil {
		return m, loadResources
	}
	return m, nil
}

func (m *Model) updateSuggestions() {
	if m.CommandInput == "" {
		m.Suggestions = nil
		return
	}

	parts := strings.Fields(m.CommandInput)
	if len(parts) == 0 {
		m.Suggestions = command.Suggest(m.CommandInput)
		return
	}

	if len(parts) == 1 && !strings.HasSuffix(m.CommandInput, " ") {
		// Still typing the command
		m.Suggestions = command.Suggest(parts[0])
	} else {
		// Typing arguments
		cmdName := parts[0]
		argQuery := ""
		if strings.HasSuffix(m.CommandInput, " ") {
			// Ready for next arg
		} else {
			argQuery = parts[len(parts)-1]
		}
		m.Suggestions = command.SuggestArgs(cmdName, argQuery, m.Containers, m.Images, m.Volumes, m.Networks)
	}
}