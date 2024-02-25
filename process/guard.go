package process

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path"
	"strings"
	"time"
)

func (p *Process) guard() {
	for {
		opts := strings.Fields(p.Options)
		cmd := exec.Command(p.GamePath, opts...)
		cmd.Dir = path.Dir(p.GamePath)
		if err := cmd.Start(); err != nil {
			log.WithField("err", err).Warning("Start game process failed")
			time.Sleep(1 * time.Second)
			continue
		}
		p.cmd = cmd
		log.Info("Game process is started.")
		cmd.Wait()
		p.lock.RLock()
		r := p.running
		p.lock.RUnlock()
		if !r {
			return
		}
		log.Warning("Game process is downed, restart...")
		time.Sleep(1 * time.Second)
	}
}
