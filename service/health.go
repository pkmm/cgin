package service

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
)

type healthService struct {
}

var HealthService healthService

const (
	//B one byte
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func (h *healthService) MemoryUseInfo() string {
	info, _ := mem.VirtualMemory()
	return fmt.Sprintf("Total: %dMB, USED: %dMB, PERCENT: %.2f%%", info.Total/MB, info.Used/MB, info.UsedPercent)
}
