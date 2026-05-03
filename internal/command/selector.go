package command

import (
	"strings"
	"uldocker/pkg/types"
)

func MatchContainers(query string, containers []types.Container) []types.Container {
	var result []types.Container

	query = strings.ToLower(query)

	for _, c := range containers {
		if strings.Contains(strings.ToLower(c.Name), query) {
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

func ResolveTargets(args []string, containers []types.Container) []types.Container {
	if len(args) == 0 {
		return nil
	}

	// First try exact keyword matching
	keywordMatch := FilterContainers(args[0], containers)
	if keywordMatch != nil {
		return keywordMatch
	}

	// Fallback to fuzzy matching
	return MatchContainers(args[0], containers)
}
