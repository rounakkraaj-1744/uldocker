package command

import (
	"strings"
	"uldocker/pkg/types"
)

var Commands = []string{
	"stop",
	"start",
	"restart",
	"rm",
	"rmi",
	"rmv",
	"rmn",
	"logs",
	"stats",
	"prune",
}

func Suggest(input string) []string {
	var res []string
	if input == "" {
		return nil
	}

	for _, cmd := range Commands {
		if strings.HasPrefix(cmd, strings.ToLower(input)) {
			res = append(res, cmd)
		}
	}

	return res
}

func SuggestArgs(cmdName string, argQuery string, containers []types.Container, images []types.Image, volumes []types.Volume, networks []types.Network) []string {
	var res []string
	
	switch cmdName {
	case "stop", "start", "restart", "rm", "logs", "stats":
		matches := MatchContainers(argQuery, containers)
		for _, c := range matches {
			res = append(res, c.Name)
		}
	case "rmi":
		matches := MatchImages(argQuery, images)
		for _, img := range matches {
			res = append(res, img.Repository)
		}
	case "rmv":
		matches := MatchVolumes(argQuery, volumes)
		for _, v := range matches {
			res = append(res, v.Name)
		}
	case "rmn":
		matches := MatchNetworks(argQuery, networks)
		for _, n := range matches {
			res = append(res, n.Name)
		}
	}

	return res
}