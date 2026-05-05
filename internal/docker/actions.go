package docker

import (
	"fmt"
	"io"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func RemoveImage(id string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	_, err = cli.ImageRemove(GetContext(), id, image.RemoveOptions{Force: true})
	return err
}

func RemoveVolume(name string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	return cli.VolumeRemove(GetContext(), name, true)
}

func RemoveNetwork(id string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	return cli.NetworkRemove(GetContext(), id)
}

func PruneContainers() (uint64, error) {
	cli, err := NewClient()
	if err != nil {
		return 0, err
	}
	report, err := cli.ContainersPrune(GetContext(), filters.Args{})
	return report.SpaceReclaimed, err
}

func PruneImages() (uint64, error) {
	cli, err := NewClient()
	if err != nil {
		return 0, err
	}
	report, err := cli.ImagesPrune(GetContext(), filters.Args{})
	return report.SpaceReclaimed, err
}

func PruneVolumes() (uint64, error) {
	cli, err := NewClient()
	if err != nil {
		return 0, err
	}
	report, err := cli.VolumesPrune(GetContext(), filters.Args{})
	return report.SpaceReclaimed, err
}

func PruneNetworks() (int, error) {
	cli, err := NewClient()
	if err != nil {
		return 0, err
	}
	report, err := cli.NetworksPrune(GetContext(), filters.Args{})
	return len(report.NetworksDeleted), err
}

func PullImage(name string) (io.ReadCloser, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}
	return cli.ImagePull(GetContext(), name, image.PullOptions{})
}

func PushImage(name string) (io.ReadCloser, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}
	return cli.ImagePush(GetContext(), name, image.PushOptions{})
}

func InspectResource(id string, resourceType string) (interface{}, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}
	ctx := GetContext()

	switch resourceType {
	case "container":
		return cli.ContainerInspect(ctx, id)
	case "image":
		raw, _, err := cli.ImageInspectWithRaw(ctx, id)
		return raw, err
	case "volume":
		return cli.VolumeInspect(ctx, id)
	case "network":
		return cli.NetworkInspect(ctx, id, network.InspectOptions{})
	}
	return nil, fmt.Errorf("invalid resource type: %s", resourceType)
}