package process

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"palworld-guard/common/rest"
	"palworld-guard/config"
	"palworld-guard/cron"
	"sync"
)

type Process struct {
	cmd     *exec.Cmd
	running bool
	lock    sync.RWMutex
	api     *rest.Rest
	*config.ProcessConfig
}

func New(c *config.ProcessConfig, api *rest.Rest) *Process {
	return &Process{
		ProcessConfig: c,
		api:           api,
		running:       false,
	}
}

func (p *Process) Start() error {
	err := cron.AddTask("memory", "*/5 * * * *", p.memoryChecker)
	if err != nil {
		return err
	}
	if len(p.PeriodicRestartInterval) != 0 {
		err = cron.AddTask("reboot", p.PeriodicRestartInterval, func() {
			err := p.SoftShutdown(60, p.MaintenanceWarningMessage)
			if err != nil {
				log.WithField("err", err).Error("Soft shutdown failed")
			}
		})
		if err != nil {
			return err
		}
	}
	if len(p.AutoBackupInterval) != 0 {
		err = cron.AddTask("backup", p.AutoBackupInterval, p.backup)
		if err != nil {
			return err
		}
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.running = true
	go p.guard()
	return nil
}

func (p *Process) SoftShutdown(sec int, msg ...string) error {
	for _, s := range msg {
		_ = p.api.AnnounceMessage(s)
	}
	ps, err := p.api.GetPlayerList()
	if err != nil {
		log.WithField("err", err).Warning("List players failed")
	} else {
		for _, player := range ps.Players {
			err := p.api.KickPlayer(player.Userid, "The server will be reboot")
			if err != nil {
				log.WithFields(map[string]interface{}{
					"player": player.Name,
					"err":    err,
				}).Warn("kick player failed")
			}
		}
	}
	err = p.api.SaveWorld()
	if err != nil {
		return fmt.Errorf("save data error: %w", err)
	}
	err = p.api.ShutdownServer(sec, "Reboot_In_60_Seconds")
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	return nil
}

func (p *Process) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.running = false
	cron.DelTask("reboot")
	cron.DelTask("memory")
	cron.DelTask("backup")
	if p.cmd == nil {
		return nil
	}
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	p.cmd.Wait()
	return nil
}
