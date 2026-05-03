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

type Image struct {
	Repository string
	Tag        string
	Size       string
}

type Volume struct {
	Name       string
	Mountpoint string
}

type Network struct {
	Name   string
	Driver string
}

type Model struct {
	ActiveTab Tab
	SelectedIndexes map[Tab]int
	ShowDetails     bool
	Width  int
	Height int
	Containers []types.Container
	Images   []Image
	Volumes  []Volume
	Networks []Network
	Loading bool
	Err     error
	Logs        []string
	Streaming   bool
	CurrentID   string
	LogsCancel  context.CancelFunc
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

		Images: []Image{
			{Repository: "node", Tag: "18", Size: "120MB"},
			{Repository: "postgres", Tag: "15", Size: "300MB"},
		},
		Volumes: []Volume{
			{Name: "pg_data", Mountpoint: "/var/lib/docker/volumes/pg_data/_data"},
			{Name: "redis_data", Mountpoint: "/var/lib/docker/volumes/redis_data/_data"},
		},
		Networks: []Network{
			{Name: "bridge", Driver: "bridge"},
			{Name: "host", Driver: "host"},
		},


		Logs:      []string{},
		Streaming: false,
		CurrentID: "",
	}
}

func (m Model) Init() tea.Cmd {
	return loadContainers
}