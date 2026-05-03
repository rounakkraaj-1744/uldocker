package command

import (
	"fmt"
	"uldocker/pkg/types"
)

var globalRegistry *Registry

func init() {
	globalRegistry = NewRegistry()
	globalRegistry.Register("stop", stopHandler)
	globalRegistry.Register("start", startHandler)
	globalRegistry.Register("restart", restartHandler)
	globalRegistry.Register("rm", rmHandler)
}

func Execute(cmd Command, targets []types.Container) (string, error) {
	handler, ok := globalRegistry.Get(cmd.Name)
	if !ok {
		return "", fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(cmd.Args, Context{Targets: targets})
}

// Handlers
func stopHandler(args []string, ctx Context) (string, error) {
	return fmt.Sprintf("Stopped %d containers", len(ctx.Targets)), nil
}

func startHandler(args []string, ctx Context) (string, error) {
	return fmt.Sprintf("Started %d containers", len(ctx.Targets)), nil
}

func restartHandler(args []string, ctx Context) (string, error) {
	return fmt.Sprintf("Restarted %d containers", len(ctx.Targets)), nil
}

func rmHandler(args []string, ctx Context) (string, error) {
	return fmt.Sprintf("Removed %d containers", len(ctx.Targets)), nil
}