package discord

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func Test_serverStatusHandle(t *testing.T) {
	var msg string
	m, err := mem.VirtualMemory()
	if err != nil {
		msg = "失敗しました"
		log.WithField("err", err).Error("get memory info failed")
		return
	}
	mU := fmt.Sprintf("%.2f", m.UsedPercent) + "%"
	c, err := cpu.Percent(time.Millisecond*100, false)
	if err != nil {
		msg = "失敗しました"
		log.WithField("err", err).Error("get cpu info failed")
		return
	}
	cpuUF := 0.0
	for _, v := range c {
		cpuUF += v
	}
	if len(c) > 1 {
		cpuUF /= float64(len(c))
	}
	cpuU := fmt.Sprintf("%.2f", cpuUF) + "%"
	msg = "CPUの使用率: " + cpuU + "\nメモリの使用率: " + mU
	t.Log(msg)
}
