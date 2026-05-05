package types

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

type ContainerStats struct {
	CPUPercentage    float64
	MemoryUsage     float64
	MemoryLimit     float64
	MemoryPercentage float64
	NetIO           string
	BlockIO         string
}