package command

import (
	"fmt"
	"uldocker/internal/docker"
	"uldocker/pkg/types"
)

var globalRegistry *Registry

func init() {
	globalRegistry = NewRegistry()
	globalRegistry.Register("stop", stopHandler)
	globalRegistry.Register("start", startHandler)
	globalRegistry.Register("restart", restartHandler)
	globalRegistry.Register("rm", rmHandler)
	globalRegistry.Register("rmi", rmiHandler)
	globalRegistry.Register("rmv", rmvHandler)
	globalRegistry.Register("rmn", rmnHandler)
	globalRegistry.Register("prune", pruneHandler)
}

func Execute(cmd Command, targets []types.Container, images []types.Image, volumes []types.Volume, networks []types.Network) (string, error) {
	handler, ok := globalRegistry.Get(cmd.Name)
	if !ok {
		return "", fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(cmd.Args, Context{
		Targets:        targets,
		ImageTargets:   images,
		VolumeTargets:  volumes,
		NetworkTargets: networks,
	})
}

// Handlers
func stopHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.Targets {
		if err := docker.StopContainer(t.ID); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Stopped %d containers", len(ctx.Targets)), nil
}

func startHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.Targets {
		if err := docker.StartContainer(t.ID); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Started %d containers", len(ctx.Targets)), nil
}

func restartHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.Targets {
		if err := docker.RestartContainer(t.ID); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Restarted %d containers", len(ctx.Targets)), nil
}

func rmHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.Targets {
		if err := docker.RemoveContainer(t.ID); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Removed %d containers", len(ctx.Targets)), nil
}

func rmiHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.ImageTargets {
		if err := docker.RemoveImage(t.Repository); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Removed %d images", len(ctx.ImageTargets)), nil
}

func rmvHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.VolumeTargets {
		if err := docker.RemoveVolume(t.Name); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Removed %d volumes", len(ctx.VolumeTargets)), nil
}

func rmnHandler(args []string, ctx Context) (string, error) {
	for _, t := range ctx.NetworkTargets {
		if err := docker.RemoveNetwork(t.Name); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("Removed %d networks", len(ctx.NetworkTargets)), nil
}

func pruneHandler(args []string, ctx Context) (string, error) {
	reclaimedContainers, _ := docker.PruneContainers()
	reclaimedImages, _ := docker.PruneImages()
	reclaimedVolumes, _ := docker.PruneVolumes()
	reclaimedNetworks, _ := docker.PruneNetworks()

	total := reclaimedContainers + reclaimedImages + reclaimedVolumes
	return fmt.Sprintf("Pruned resources. Reclaimed %.1fMB (%d networks deleted)", float64(total)/1024.0/1024.0, reclaimedNetworks), nil
}