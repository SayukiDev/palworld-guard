package process

import (
	"context"
	arc "github.com/mholt/archiver/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const saveDataPath = "Pal/Saved/SaveGames/0"

func (p *Process) backup() {
	log.Info("Save data backup...")
	savePath := path.Join(path.Dir(p.GamePath), saveDataPath)
	_, err := os.Stat(savePath)
	if err != nil {
		log.WithField("err", err).Error("Check save data failed")
		return
	}
	fs, err := arc.FilesFromDisk(nil, map[string]string{
		savePath: "backup",
	})
	if err != nil {
		log.WithField("err", err).Error("Not found save data")
		return
	}
	filename := "Backup-" + time.Now().Format("2006-01-02|15:04")
	f, err := os.Create(path.Join(p.BackupPath, filename))
	if err != nil {
		log.WithField("err", err).Error("Create package failed")
		return
	}
	defer f.Close()
	z := arc.Zip{}
	err = z.Archive(context.Background(), f, fs)
	if err != nil {
		log.WithField("err", err).Error("Write package failed")
		return
	}
	log.WithField("output", filename).Info("Save data backup done")
}
