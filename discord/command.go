package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (d *Discord) initCommand() error {
	cmds := []discordgo.ApplicationCommand{
		{
			Name:        "restart",
			Description: "パルワールドのサーバーを再起動します",
		},
		{
			Name:        "list-players",
			Description: "ログインしてるプレイヤーの情報を取得します",
		},
		{
			Name:        "server-status",
			Description: "サーバーの状態をチェックします",
		},
	}
	d.commandHandles = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"restart":       d.restartHandle,
		"list-players":  d.listPlayersHandle,
		"server-status": d.serverStatusHandle,
	}
	for i := range cmds {
		_, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, "", &cmds[i])
		if err != nil {
			return fmt.Errorf("add command %s error: %s", cmds[i].Name, err)
		}
	}
	d.s.AddHandler(d.handleSelector)
	return nil
}

func (d *Discord) handleSelector(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		d.commandHandles[i.ApplicationCommandData().Name](s, i)
	}
}
