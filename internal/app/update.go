package app

import (
	"context"
	"dawker/internal/command"
	"dawker/internal/docker"
	"dawker/pkg/types"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/pkg/stdcopy"
)

type resourcesLoadedMsg struct {
	containers []types.Container
	images     []types.Image
	volumes    []types.Volume
	networks   []types.Network
	err        error
}

type logStreamStartedMsg struct {
	reader io.Reader
}

type logChunkMsg struct {
	reader io.Reader
	chunk  string
}

type streamEndedMsg struct{}

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
		m.Err = nil
		m.Containers = msg.containers
		m.Images = msg.images
		m.Volumes = msg.volumes
		m.Networks = msg.networks
		m.clampAllIndexes()
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case logStreamStartedMsg:
		return m, readStreamCmd(msg.reader)

	case logChunkMsg:
		return m.handleLogChunk(msg)

	case streamEndedMsg:
		m.Streaming = false
		return m, loadResources
	}

	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.CommandResult != "" {
		m.CommandResult = ""
	}
	if m.CommandError != "" && !m.CommandMode {
		m.CommandError = ""
	}

	if m.CommandMode {
		return m.handleCommandModeKey(msg)
	}

	return m.handleNormalModeKey(msg)
}

func (m Model) handleCommandModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return m.executeCommand()

	case tea.KeyEsc:
		m.CommandMode = false
		m.CommandInput = ""
		m.CommandError = ""
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

	case tea.KeyTab:
		if len(m.Suggestions) > 0 {
			parts := strings.Fields(m.CommandInput)
			if len(parts) <= 1 && !strings.HasSuffix(m.CommandInput, " ") {
				m.CommandInput = m.Suggestions[0] + " "
			} else {
				if len(parts) > 1 {
					parts[len(parts)-1] = m.Suggestions[0]
				} else {
					parts = append(parts, m.Suggestions[0])
				}
				m.CommandInput = strings.Join(parts, " ") + " "
			}
			m.updateSuggestions()
		}
		return m, nil

	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.CommandInput) > 0 {
			m.CommandInput = m.CommandInput[:len(m.CommandInput)-1]
		}
		m.updateSuggestions()
		return m, nil

	case tea.KeySpace:
		m.CommandInput += " "
		m.updateSuggestions()
		return m, nil

	default:
		ch := msg.String()
		if len(ch) == 1 || (len(ch) > 1 && ch[0] != 0x1b) {
			m.CommandInput += ch
			m.updateSuggestions()
		}
		return m, nil
	}
}

func (m Model) handleNormalModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "1", "2", "3", "4", "j", "k", "up", "down", "esc":
		m.cancelStream()
	}

	switch key {
	case "q", "ctrl+c":
		m.cancelStream()
		return m, tea.Quit

	case ":":
		m.CommandMode = true
		m.CommandInput = ""
		m.CommandError = ""
		m.CommandResult = ""
		m.Suggestions = nil
		m.HistoryIndex = -1
		return m, nil

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
			if idx < len(m.Containers) {
				container := m.Containers[idx]

				m.cancelStream()

				m.Logs = []string{}
				m.Streaming = true
				m.CurrentID = container.ID

				ctx, cancel := context.WithCancel(context.Background())
				m.LogsCancel = cancel

				return m, startStreamCmd(ctx, container.ID)
			}
		}
		m.ShowDetails = true
	}

	return m, nil
}

func (m Model) handleLogChunk(msg logChunkMsg) (tea.Model, tea.Cmd) {
	lines := strings.Split(msg.chunk, "\n")

	if len(m.Logs) == 0 {
		m.Logs = lines
	} else {
		m.Logs[len(m.Logs)-1] += lines[0]
		if len(lines) > 1 {
			m.Logs = append(m.Logs, lines[1:]...)
		}
	}

	if len(m.Logs) > 500 {
		m.Logs = m.Logs[len(m.Logs)-500:]
	}

	return m, readStreamCmd(msg.reader)
}

func (m *Model) moveDown() {
	max := m.getActiveListLength()
	if max == 0 {
		return
	}
	idx := m.SelectedIndexes[m.ActiveTab]
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

func (m *Model) clampAllIndexes() {
	for _, tab := range []Tab{TabContainers, TabImages, TabVolumes, TabNetworks} {
		max := 0
		switch tab {
		case TabContainers:
			max = len(m.Containers)
		case TabImages:
			max = len(m.Images)
		case TabVolumes:
			max = len(m.Volumes)
		case TabNetworks:
			max = len(m.Networks)
		}
		if max == 0 {
			m.SelectedIndexes[tab] = 0
		} else if m.SelectedIndexes[tab] >= max {
			m.SelectedIndexes[tab] = max - 1
		}
	}
}

func (m *Model) cancelStream() {
	if m.LogsCancel != nil {
		m.LogsCancel()
		m.LogsCancel = nil
	}
	m.Streaming = false
}

func startPullCmd(ctx context.Context, imageName string) tea.Cmd {
	return func() tea.Msg {
		reader, err := docker.PullImage(imageName)
		if err != nil {
			return logChunkMsg{chunk: "ERROR: " + err.Error()}
		}
		return logStreamStartedMsg{reader: reader}
	}
}

func startPushCmd(ctx context.Context, imageName string) tea.Cmd {
	return func() tea.Msg {
		reader, err := docker.PushImage(imageName)
		if err != nil {
			return logChunkMsg{chunk: "ERROR: " + err.Error()}
		}
		return logStreamStartedMsg{reader: reader}
	}
}

func (m Model) executeCommand() (tea.Model, tea.Cmd) {
	cmd := command.Parse(m.CommandInput)
	if cmd.Name == "" {
		m.CommandMode = false
		return m, nil
	}

	var targets []types.Container
	var images []types.Image
	var volumes []types.Volume
	var networks []types.Network

	if len(cmd.Args) == 0 {
		idx := m.SelectedIndexes[m.ActiveTab]
		switch m.ActiveTab {
		case TabContainers:
			if idx < len(m.Containers) {
				targets = []types.Container{m.Containers[idx]}
			}
		case TabImages:
			if idx < len(m.Images) {
				images = []types.Image{m.Images[idx]}
			}
		case TabVolumes:
			if idx < len(m.Volumes) {
				volumes = []types.Volume{m.Volumes[idx]}
			}
		case TabNetworks:
			if idx < len(m.Networks) {
				networks = []types.Network{m.Networks[idx]}
			}
		}
	} else {
		switch cmd.Name {
		case "rmi", "pull":
			images = command.MatchImages(cmd.Args[0], m.Images)
		case "rmv":
			volumes = command.MatchVolumes(cmd.Args[0], m.Volumes)
		case "rmn":
			networks = command.MatchNetworks(cmd.Args[0], m.Networks)
		case "prune", "inspect":
			// handled inside execute
		default:
			targets = command.ResolveTargets(cmd.Args, m.Containers)
		}
	}

	resultMsg, err := command.Execute(cmd, targets, images, volumes, networks)

	m.CommandMode = false
	m.Suggestions = nil

	if err != nil {
		m.CommandError = err.Error()
		m.CommandResult = ""
		return m, nil
	}

	// Handle special prefixes
	if strings.HasPrefix(resultMsg, "PULL:") {
		imageName := strings.TrimPrefix(resultMsg, "PULL:")
		m.cancelStream()
		m.Logs = []string{"Pulling " + imageName + "..."}
		m.Streaming = true
		ctx, cancel := context.WithCancel(context.Background())
		m.LogsCancel = cancel
		return m, startPullCmd(ctx, imageName)
	}

	if strings.HasPrefix(resultMsg, "PUSH:") {
		imageName := strings.TrimPrefix(resultMsg, "PUSH:")
		m.cancelStream()
		m.Logs = []string{"Pushing " + imageName + "..."}
		m.Streaming = true
		ctx, cancel := context.WithCancel(context.Background())
		m.LogsCancel = cancel
		return m, startPushCmd(ctx, imageName)
	}

	if strings.HasPrefix(resultMsg, "STATS:") {
		// For now, just show a message. Proper stats implementation needs a loop.
		m.CommandResult = "Stats view coming soon"
		return m, nil
	}

	if cmd.Name == "inspect" {
		m.cancelStream()
		m.Logs = strings.Split(resultMsg, "\n")
		m.Streaming = true
		m.ShowDetails = true
		return m, nil
	}

	m.CommandError = ""
	m.CommandResult = resultMsg

	m.History = append(m.History, m.CommandInput)
	if len(m.History) > 50 {
		m.History = m.History[len(m.History)-50:]
	}
	m.HistoryIndex = -1

	return m, loadResources
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
		m.Suggestions = command.Suggest(parts[0])
	} else {
		cmdName := parts[0]
		argQuery := ""
		if !strings.HasSuffix(m.CommandInput, " ") {
			argQuery = parts[len(parts)-1]
		}
		m.Suggestions = command.SuggestArgs(cmdName, argQuery, m.Containers, m.Images, m.Volumes, m.Networks)
	}
}
