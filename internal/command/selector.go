package command

import (
	"dawker/pkg/types"
	"strings"
)

func MatchContainers(query string, containers []types.Container) []types.Container {
	var result []types.Container

	query = strings.ToLower(query)

	for _, c := range containers {
		if strings.Contains(strings.ToLower(c.Name), query) || strings.Contains(strings.ToLower(c.ID), query) {
			result = append(result, c)
		}
	}

	return result
}

func filterByState(containers []types.Container, state string) []types.Container {
	var result []types.Container
	for _, c := range containers {
		if strings.Contains(strings.ToLower(c.Status), state) {
			result = append(result, c)
		}
	}
	return result
}

func FilterContainers(cmd string, containers []types.Container) []types.Container {
	switch cmd {
	case "all":
		return containers

	case "running":
		return filterByState(containers, "up")

	case "exited":
		return filterByState(containers, "exited")

	default:
		return nil
	}
}

func MatchImages(query string, images []types.Image) []types.Image {
	var result []types.Image
	query = strings.ToLower(query)
	for _, i := range images {
		if strings.Contains(strings.ToLower(i.Repository), query) || strings.Contains(strings.ToLower(i.ID), query) {
			result = append(result, i)
		}
	}
	return result
}

func MatchVolumes(query string, volumes []types.Volume) []types.Volume {
	var result []types.Volume
	query = strings.ToLower(query)
	for _, v := range volumes {
		if strings.Contains(strings.ToLower(v.Name), query) {
			result = append(result, v)
		}
	}
	return result
}

func MatchNetworks(query string, networks []types.Network) []types.Network {
	var result []types.Network
	query = strings.ToLower(query)
	for _, n := range networks {
		if strings.Contains(strings.ToLower(n.Name), query) {
			result = append(result, n)
		}
	}
	return result
}

func ResolveTargets(args []string, containers []types.Container) []types.Container {
	if len(args) == 0 {
		return nil
	}

	keywordMatch := FilterContainers(args[0], containers)
	if keywordMatch != nil {
		return keywordMatch
	}

	return MatchContainers(args[0], containers)
}