package app

import (
	"context"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/pkg/stdcopy"
	"uldocker/internal/docker"
	"uldocker/pkg/types"
)

type containersLoadedMsg struct {
	containers []types.Container
	err        error
}

func loadContainers() tea.Msg {
	containers, err := docker.ListContainers()
	return containersLoadedMsg{
		containers: containers,
		err:        err,
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

	case containersLoadedMsg:
		m.Loading = false
		m.Containers = msg.containers
		m.Err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "1", "2", "3", "4", "j", "k", "up", "down":
			if msg.String() != "enter" && m.LogsCancel != nil {
				m.LogsCancel()
				m.LogsCancel = nil
				m.Streaming = false
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			m.Loading = true
			return m, loadContainers

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