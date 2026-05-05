package docker

import (
	"dawker/pkg/types"
	"encoding/json"
	"fmt"
	"io"

	container "github.com/docker/docker/api/types/container"
)

func ParseStats(reader io.Reader) (types.ContainerStats, error) {
	var s container.StatsResponse
	if err := json.NewDecoder(reader).Decode(&s); err != nil {
		return types.ContainerStats{}, err
	}

	cpuPercent := calculateCPUPercentUnix(&s)
	memUsage := float64(s.MemoryStats.Usage)
	memLimit := float64(s.MemoryStats.Limit)
	memPercent := 0.0
	if memLimit > 0 {
		memPercent = (memUsage / memLimit) * 100.0
	}

	return types.ContainerStats{
		CPUPercentage:    cpuPercent,
		MemoryUsage:     memUsage,
		MemoryLimit:     memLimit,
		MemoryPercentage: memPercent,
		NetIO:           fmt.Sprintf("%.1fMB / %.1fMB", float64(s.Networks["eth0"].RxBytes)/1024/1024, float64(s.Networks["eth0"].TxBytes)/1024/1024),
		BlockIO:         fmt.Sprintf("%.1fMB / %.1fMB", float64(s.StorageStats.ReadCountNormalized)/1024/1024, float64(s.StorageStats.WriteCountNormalized)/1024/1024),
	}, nil
}

func calculateCPUPercentUnix(v *container.StatsResponse) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
		onlineCPUs  = float64(v.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}
