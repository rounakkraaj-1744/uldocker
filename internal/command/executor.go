package command

import (
	"dawker/internal/docker"
	"dawker/pkg/types"
	"encoding/json"
	"fmt"
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
	globalRegistry.Register("pull", pullHandler)
	globalRegistry.Register("inspect", inspectHandler)
	globalRegistry.Register("stats", statsHandler)
	globalRegistry.Register("run", runHandler)
	globalRegistry.Register("push", pushHandler)
	globalRegistry.Register("exec", execHandler)
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

func pullHandler(args []string, ctx Context) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("image name required")
	}
	return "PULL:" + args[0], nil
}

func inspectHandler(args []string, ctx Context) (string, error) {
	var targetID string
	var resourceType string

	if len(ctx.Targets) > 0 {
		targetID = ctx.Targets[0].ID
		resourceType = "container"
	} else if len(ctx.ImageTargets) > 0 {
		targetID = ctx.ImageTargets[0].Repository
		resourceType = "image"
	} else if len(ctx.VolumeTargets) > 0 {
		targetID = ctx.VolumeTargets[0].Name
		resourceType = "volume"
	} else if len(ctx.NetworkTargets) > 0 {
		targetID = ctx.NetworkTargets[0].Name
		resourceType = "network"
	}

	if targetID == "" {
		return "", fmt.Errorf("no target selected for inspect")
	}

	data, err := docker.InspectResource(targetID, resourceType)
	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func statsHandler(args []string, ctx Context) (string, error) {
	if len(ctx.Targets) == 0 {
		return "", fmt.Errorf("no container selected for stats")
	}
	return "STATS:" + ctx.Targets[0].ID, nil
}

func runHandler(args []string, ctx Context) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("image name required")
	}
	id, err := docker.CreateAndStartContainer(args[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Started container %s", id), nil
}

func pushHandler(args []string, ctx Context) (string, error) {
	if len(args) == 0 && len(ctx.ImageTargets) == 0 {
		return "", fmt.Errorf("image name required")
	}
	var name string
	if len(args) > 0 {
		name = args[0]
	} else {
		name = ctx.ImageTargets[0].Repository
	}
	return "PUSH:" + name, nil
}

func execHandler(args []string, ctx Context) (string, error) {
	if len(ctx.Targets) == 0 {
		return "", fmt.Errorf("no container selected for exec")
	}
	shell := "/bin/sh"
	if len(args) > 0 {
		shell = args[0]
	}
	return "EXEC:" + ctx.Targets[0].ID + ":" + shell, nil
}

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