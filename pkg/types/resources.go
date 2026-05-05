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