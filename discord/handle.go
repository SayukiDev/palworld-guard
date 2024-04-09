package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
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
	a := d.api
	_ = a.AnnounceMessage("The_Server_Will_Reboot")
	ps, err := a.GetPlayerList()
	if err != nil {
		log.WithField("err", err).Warning("List players failed")
	} else {
		for _, player := range ps.Players {
			a.KickPlayer(player.Userid, "The server will be reboot")
		}
	}
	a.SaveWorld()
	a.ShutdownServer(10, "Reboot_In_10_Seconds")
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
	a := d.api
	ps, err := a.GetPlayerList()
	if err != nil {
		log.WithField("err", err).Warning("List players failed")
		msg = "失敗しました"
		return
	}
	msg = fmt.Sprintf("今ログインしてるプレイヤーは%d人です", len(ps.Players))
	if len(ps.Players) > 0 {
		msg += "\n\n----------------\n"
		for _, p := range ps.Players {
			msg += fmt.Sprintf("名前: %s Id: %s\n座標: X %.2f Y%.2f\nレベル: %d\n",
				p.Name, p.Userid,
				p.LocationX, p.LocationY,
				p.Level)
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
