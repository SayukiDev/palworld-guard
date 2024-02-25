package main

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"os/signal"
	"palworld-guard/config"
	"palworld-guard/cron"
	"palworld-guard/discord"
	"palworld-guard/process"
	"syscall"
	"time"
)

func main() {

	log.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.DateTime,
	})
	log.Info("Starting...")
	c := config.New("./config.json5")
	log.Info("Config Loading...")
	err := c.Load()
	if err != nil {
		log.WithField("err", err).
			Warning("Load config file failed, use default config")
	} else {
		log.Info("Config loaded.")
	}
	switch c.LogLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
	log.Info("Staring guard...")
	p := process.New(c.Process, c.Rcon)
	_ = p.Start()
	log.Info("Guard started.")
	if c.Discord.Enable {
		log.Info("Starting Discord...")
		d, err := discord.New(c.Discord, c.Rcon)
		if err != nil {
			log.WithField("err", err).Error("Init Discord failed")
			return
		}
		err = d.Start()
		if err != nil {
			log.WithField("err", err).Error("Start Discord failed")
			return
		}
		log.Info("Discord started.")
	}
	_ = cron.Start()
	log.Info("Done. Press CTRL-C to exit.")

	waitExit()
	// Exit process
	log.Info("Exiting...")
	_ = cron.Close()
	_ = p.Close()
	log.Info("Exited.")
}

func waitExit() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sc
}
