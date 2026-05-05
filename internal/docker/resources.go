package docker

import (
	"dawker/pkg/types"
	"fmt"
	"strings"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

func ListImages() ([]types.Image, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}

	ctx := GetContext()

	images, err := cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, err
	}

	var result []types.Image
	for _, img := range images {
		repo := "<none>"
		tag := "<none>"

		if len(img.RepoTags) > 0 {
			parts := strings.Split(img.RepoTags[0], ":")
			if len(parts) == 2 {
				repo = parts[0]
				tag = parts[1]
			} else {
				repo = img.RepoTags[0]
			}
		}

		sizeMB := float64(img.Size) / 1024.0 / 1024.0
		sizeStr := fmt.Sprintf("%.1fMB", sizeMB)

		id := img.ID
		if strings.HasPrefix(id, "sha256:") {
			if len(id) >= 19 {
				id = id[7:19]
			} else {
				id = id[7:]
			}
		} else if len(id) > 12 {
			id = id[:12]
		}

		result = append(result, types.Image{
			ID:         id,
			Repository: repo,
			Tag:        tag,
			Size:       sizeStr,
		})
	}
	return result, nil
}

func ListVolumes() ([]types.Volume, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}

	ctx := GetContext()

	vols, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []types.Volume
	for _, v := range vols.Volumes {
		result = append(result, types.Volume{
			Name:       v.Name,
			Mountpoint: v.Mountpoint,
		})
	}
	return result, nil
}

func ListNetworks() ([]types.Network, error) {
	cli, err := NewClient()
	if err != nil {
		return nil, err
	}

	ctx := GetContext()

	nets, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []types.Network
	for _, n := range nets {
		result = append(result, types.Network{
			Name:   n.Name,
			Driver: n.Driver,
		})
	}
	return result, nil
}