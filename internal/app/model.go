package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabContainers Tab = iota
	TabImages
	TabVolumes
	TabNetworks
)

// Mock data structures
type Container struct {
	Name   string
	ID     string
	Status string // "running" or "exited"
}

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

// Model represents the UI state
type Model struct {
	ActiveTab    Tab
	SelectedIndexes map[Tab]int
	ShowDetails  bool

	Width  int
	Height int

	// Mock Data
	Containers []Container
	Images     []Image
	Volumes    []Volume
	Networks   []Network
}

func NewModel() Model {
	return Model{
		ActiveTab:       TabContainers,
		SelectedIndexes: map[Tab]int{
			TabContainers: 0,
			TabImages:     0,
			TabVolumes:    0,
			TabNetworks:   0,
		},
		ShowDetails: false,

		Containers: []Container{
			{Name: "api-server", ID: "a1b2c3d4e5f6", Status: "running"},
			{Name: "postgres-db", ID: "b2c3d4e5f6a1", Status: "running"},
			{Name: "redis-cache", ID: "c3d4e5f6a1b2", Status: "exited"},
		},
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
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}