package config

import (
	"encoding/json"
	"os"
	"palworld-guard/common/jsontrim"
	"sync"
)

type Config struct {
	path     string
	lock     sync.Mutex
	LogLevel string         `json:"LogLevel"`
	Process  *ProcessConfig `json:"Process"`
	Rcon     *RconConfig    `json:"Rcon"`
	Discord  *DiscordConfig `json:"Discord"`
}

type ProcessConfig struct {
	GamePath                  string  `json:"GamePath"`
	MemoryUsageThreshold      float64 `json:"MemoryUsageThreshold"`
	MaintenanceWarningMessage string  `json:"MaintenanceWarningMessage"`
	PeriodicRestartInterval   string  `json:"PeriodicRestartInterval"` // crontab format
	AutoBackupInterval        string  `json:"AutoBackupInterval"`
	BackupPath                string  `json:"BackupPath"`
	Options                   string  `json:"StartOptions"`
}

type RconConfig struct {
	Addr          string `json:"Address"`
	AdminPassword string `json:"AdminPassword"`
}

type DiscordConfig struct {
	Enable  bool     `json:"Enable"`
	Token   string   `json:"Token"`
	Masters []string `json:"Masters"`
}

func New(p string) *Config {
	return &Config{
		LogLevel: "info",
		Rcon: &RconConfig{
			Addr:          "127.0.0.1:25575",
			AdminPassword: "",
		},
		Process: &ProcessConfig{
			GamePath:                  "/home/pal/Steam/steamapps/common/PalServer/PalServer.sh",
			MemoryUsageThreshold:      97,
			MaintenanceWarningMessage: "Memory_Not_Enough_The_Server_Will_Reboot",
			PeriodicRestartInterval:   "0 6 * * *",
			AutoBackupInterval:        "",
			BackupPath:                "/home/pal/backup/",
			Options:                   "EpicApp=PalServer -useperfthreads -NoAsyncLoadingThread -UseMultithreadForDS",
		},
		Discord: &DiscordConfig{
			Enable:  false,
			Token:   "token",
			Masters: make([]string, 0),
		},
		path: p,
	}
}

func (c *Config) Load() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	f, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(jsontrim.NewTrimNodeReader(f)).Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Save() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		return err
	}
	return nil
}
