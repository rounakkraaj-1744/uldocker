package docker

import (
	"context"
	"github.com/docker/docker/client"
)

func NewClient () (*client.Client, error){
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}

func GetContext() context.Context {
	return context.Background()
}