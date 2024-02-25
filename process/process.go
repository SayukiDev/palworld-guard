package process

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"palworld-guard/common/rcon"
	"palworld-guard/config"
	"palworld-guard/cron"
	"sync"
)

type Process struct {
	cmd     *exec.Cmd
	running bool
	lock    sync.RWMutex
	rconC   *config.RconConfig
	*config.ProcessConfig
}

func New(c *config.ProcessConfig, rconC *config.RconConfig) *Process {
	return &Process{
		ProcessConfig: c,
		rconC:         rconC,
		running:       false,
	}
}

func (p *Process) Start() error {
	err := cron.AddTask("memory", "*/5 * * * *", p.memoryChecker)
	if err != nil {
		return err
	}
	err = cron.AddTask("periodic", p.PeriodicRestartInterval, func() {
		err := p.SoftShutdown("60", p.MaintenanceWarningMessage)
		if err != nil {
			log.WithField("err", err).Error("Soft shutdown failed")
		}
	})
	if err != nil {
		return err
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.running = true
	go p.guard()
	return nil
}

func (p *Process) SoftShutdown(sec string, msg ...string) error {
	rcon, err := rcon.NewPalRcon(p.rconC.Addr, p.rconC.AdminPassword)
	if err != nil {
		return fmt.Errorf("connect to rcon error: %w", err)
	}
	defer rcon.Close()
	for _, s := range msg {
		rcon.Broadcast(s)
	}
	ps, err := rcon.ShowPlayers()
	if err != nil {
		log.WithField("err", err).Warning("List players failed")
	} else {
		for _, player := range ps {
			err := rcon.KickPlayer(player.Id)
			if err != nil {
				log.WithFields(map[string]interface{}{
					"player": player.Name,
					"err":    err,
				}).Warn("kick player failed")
			}
		}
	}
	err = rcon.Save()
	if err != nil {
		return fmt.Errorf("save data error: %w", err)
	}
	err = rcon.Shutdown(sec, "Reboot_In_60_Seconds")
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	return nil
}

func (p *Process) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.running = false
	if p.cmd == nil {
		return nil
	}
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	p.cmd.Wait()
	return nil
}
