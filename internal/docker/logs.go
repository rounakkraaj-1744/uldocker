package docker

import (
	"context"
	"io"
	"github.com/docker/docker/api/types/container"
)

func StreamLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}

	return cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "50",
	})
}