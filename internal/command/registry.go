package command

import (
	"uldocker/pkg/types"
)

type Context struct {
	Targets []types.Container
}

type Handler func(args []string, ctx Context) (string, error)

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

func (r *Registry) Register(name string, handler Handler) {
	r.handlers[name] = handler
}

func (r *Registry) Get(name string) (Handler, bool) {
	h, ok := r.handlers[name]
	return h, ok
}
