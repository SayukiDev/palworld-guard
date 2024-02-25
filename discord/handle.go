package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
	"palworld-guard/common/rcon"
	"time"
)

func (d *Discord) pingHandle(_ *discordgo.Session, m *discordgo.InteractionCreate) {
	err := d.interactionRespondText(m, "Pong!", nil)
	if err != nil {
		log.Error("reply message error: ", err)
		return
	}
}

func (d *Discord) restartHandle(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	if !d.isMaster(i) {
		err := d.interactionRespondText(i, "権限不足", nil)
		if err != nil {
			log.Error("reply message error: ", err)
			return
		}
		return
	}
	var err error
	var msg string
	err = d.interactionRespondText(i, "実行中...", nil)
	if err != nil {
		log.Error("reply message error: ", err)
		return
	}
	defer func() {
		err = d.editInteractionResponseText(i, msg)
		if err != nil {
			log.Error("reply message error: ", err)
			return
		}
	}()
	rcon, err := rcon.NewPalRcon(d.rconC.Addr, d.rconC.AdminPassword)
	if err != nil {
		log.WithField("err", err).Error("connect to rcon failed")
		msg = "失敗しました"
		return
	}
	defer rcon.Close()
	rcon.Broadcast("The_Server_Will_Reboot")
	ps, err := rcon.ShowPlayers()
	if err != nil {
		log.WithField("err", err).Warning("List players failed")
	} else {
		for _, player := range ps {
			rcon.KickPlayer(player.Id)
		}
	}
	rcon.Save()
	rcon.Shutdown("10", "Reboot_In_10_Seconds")
	msg = "サーバーが十秒後に再起動します"
}

func (d *Discord) listPlayersHandle(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	if !d.isMaster(i) {
		err := d.interactionRespondText(i, "権限不足", nil)
		if err != nil {
			log.Error("reply message error: ", err)
			return
		}
		return
	}
	var err error
	var msg string
	err = d.interactionRespondText(i, "実行中...", nil)
	if err != nil {
		log.Error("reply message error: ", err)
		return
	}
	defer func() {
		err = d.editInteractionResponseText(i, msg)
		if err != nil {
			log.Error("reply message error: ", err)
			return
		}
	}()
	rcon, err := rcon.NewPalRcon(d.rconC.Addr, d.rconC.AdminPassword)
	if err != nil {
		log.WithField("err", err).Error("Connect to rcon failed")
		msg = "失敗しました"
		return
	}
	defer rcon.Close()
	ps, err := rcon.ShowPlayers()
	if err != nil {
		log.WithField("err", err).Warning("List players failed")
		msg = "失敗しました"
		return
	}
	msg = fmt.Sprintf("今ログインしてるプレイヤーは%d人です", len(ps))
	if len(ps) > 0 {
		msg += "\n\n----------------\n"
		for _, p := range ps {
			msg += fmt.Sprintf("名前: %s Id: %s\n", p.Name, p.Id)
		}
		msg += "----------------"
	}
}

func (d *Discord) serverStatusHandle(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	if !d.isMaster(i) {
		err := d.interactionRespondText(i, "権限不足", nil)
		if err != nil {
			log.Error("reply message error: ", err)
			return
		}
		return
	}
	var err error
	var msg string
	err = d.interactionRespondText(i, "実行中...", nil)
	if err != nil {
		log.Error("reply message error: ", err)
		return
	}
	defer func() {
		err = d.editInteractionResponseText(i, msg)
		if err != nil {
			log.Error("reply message error: ", err)
			return
		}
	}()

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
}
