package docker

import (
	"strings"
	"uldocker/pkg/types"
	"github.com/docker/docker/api/types/container"
)


func ListContainers() ([]types.Container, error) {
	cli, err:=NewClient()

	if err != nil {
		return nil, err
	}

	ctx:=GetContext()

	containers, err:=cli.ContainerList(ctx, container.ListOptions{All: true})
	if err!=nil{
		return nil, err
	}

	var result[]types.Container

	for _, c := range containers {
		name:=""
		if len (c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		result = append(result, types.Container{
			ID:     c.ID[:12],
			Name:   name,
			Status: c.Status,
		})
	}
	return result, nil
}

func StartContainer(id string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	return cli.ContainerStart(GetContext(), id, container.StartOptions{})
}

func StopContainer(id string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	return cli.ContainerStop(GetContext(), id, container.StopOptions{})
}

func RestartContainer(id string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	return cli.ContainerRestart(GetContext(), id, container.StopOptions{})
}

func RemoveContainer(id string) error {
	cli, err := NewClient()
	if err != nil {
		return err
	}
	return cli.ContainerRemove(GetContext(), id, container.RemoveOptions{Force: true})
}