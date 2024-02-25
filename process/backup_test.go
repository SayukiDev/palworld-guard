package process

import (
	"palworld-guard/config"
	"testing"
)

func TestProcess_backup(t *testing.T) {
	c := config.New("../example/config.json5")
	c.Process.BackupPath = "/tmp"
	c.Process.GamePath = "../testdata/pal_sav_test/1"
	p := New(c.Process, c.Rcon)
	p.backup()
}
