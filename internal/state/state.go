package state

import (
	"dawker/pkg/types"
)

type AppState struct {
	Containers    []types.Container
	SelectedIndex int
	IsCommandMode bool
	CommandInput  string
	LastAction    string
}

func NewAppState() *AppState {
	return &AppState{
		Containers:    []types.Container{},
		SelectedIndex: 0,
		IsCommandMode: false,
		CommandInput:  "",
	}
}