package process

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
)

func (p *Process) memoryChecker() {
	log.Info("Checking memory usage...")
	var total, used float64
	mem.VirtualMemory()
	m1, err := mem.VirtualMemory()
	if err != nil {
		log.WithField("err", err).Error("Get memory info failed")
		return
	}
	used = float64(m1.Used)
	total = float64(m1.Total)
	m2, err := mem.SwapMemory()
	if err != nil {
		log.WithField("err", err).Warning("Get swap info failed")
	} else {
		used += float64(m2.Used)
		total += float64(m2.Total)
	}
	usedP := used / total * 100
	if usedP < p.MemoryUsageThreshold {
		log.WithFields(map[string]interface{}{
			"Used":  fmt.Sprintf("%.2f", usedP) + "%",
			"Limit": fmt.Sprintf("%.2f", p.MemoryUsageThreshold) + "%",
		}).Info("Not over, pass")
		return
	}
	log.Info("Memory over, restart server after 60s")
	// Shutdown
	err = p.SoftShutdown(
		"60",
		fmt.Sprintf("Memory_Is_Above_%v%%", p.MemoryUsageThreshold),
		p.MaintenanceWarningMessage,
	)
	if err != nil {
		log.WithField("err", err).Error("Soft shutdown failed")
		return
	}
}
