package command

import (
	"dawker/pkg/types"
	"strings"
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
	"pull",
	"run",
	"inspect",
	"push",
	"exec",
}

func Suggest(input string) []string {
	var res []string
	if input == "" {
		return nil
	}

	lower := strings.ToLower(input)
	for _, cmd := range Commands {
		if strings.HasPrefix(cmd, lower) {
			res = append(res, cmd)
		}
	}

	return res
}

func SuggestArgs(cmdName string, argQuery string, containers []types.Container, images []types.Image, volumes []types.Volume, networks []types.Network) []string {
	var res []string

	switch cmdName {
	case "stop", "start", "restart", "rm", "logs", "stats":
		if argQuery == "" {
			return []string{"all", "running", "exited"}
		}
		for _, kw := range []string{"all", "running", "exited"} {
			if strings.HasPrefix(kw, strings.ToLower(argQuery)) {
				res = append(res, kw)
			}
		}
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

	if len(res) > 5 {
		res = res[:5]
	}

	return res
}

func IsValidCommand(name string) bool {
	for _, cmd := range Commands {
		if cmd == name {
			return true
		}
	}
	return false
}

func ClosestCommand(input string) string {
	input = strings.ToLower(input)
	bestMatch := ""
	bestScore := 0

	for _, cmd := range Commands {
		score := 0
		minLen := len(input)
		if len(cmd) < minLen {
			minLen = len(cmd)
		}
		for i := 0; i < minLen; i++ {
			if input[i] == cmd[i] {
				score++
			} else {
				break
			}
		}
		if score > bestScore && score >= 2 {
			bestScore = score
			bestMatch = cmd
		}
	}

	return bestMatch
}