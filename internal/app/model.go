package app

import (
	"context"
	"uldocker/pkg/types"
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabContainers Tab = iota
	TabImages
	TabVolumes
	TabNetworks
)

type Model struct {
	ActiveTab Tab
	SelectedIndexes map[Tab]int
	ShowDetails     bool
	Width  int
	Height int
	Containers []types.Container
	Images   []types.Image
	Volumes  []types.Volume
	Networks []types.Network
	Loading bool
	Err     error
	Logs        []string
	Streaming   bool
	CurrentID   string
	LogsCancel  context.CancelFunc
	CommandMode  bool
	CommandInput string
	CommandError string
	CommandResult string
	Suggestions  []string
	History      []string
	HistoryIndex int
}

func NewModel() Model {
	return Model{
		ActiveTab: TabContainers,
		SelectedIndexes: map[Tab]int{
			TabContainers: 0,
			TabImages:     0,
			TabVolumes:    0,
			TabNetworks:   0,
		},
		ShowDetails: false,

		Loading: true,

		Containers: []types.Container{},
		Images:     []types.Image{},
		Volumes:    []types.Volume{},
		Networks:   []types.Network{},


		Logs:      []string{},
		Streaming: false,
		CurrentID: "",

		Suggestions:  []string{},
		History:      []string{},
		HistoryIndex: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return loadResources
}