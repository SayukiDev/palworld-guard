package discord

import (
	dis "github.com/bwmarrin/discordgo"
	"palworld-guard/config"
)

type Discord struct {
	s              *dis.Session
	masters        []string
	rconC          *config.RconConfig
	commandHandles map[string]func(s *dis.Session, i *dis.InteractionCreate)
}

func New(c *config.DiscordConfig, rcon *config.RconConfig) (*Discord, error) {
	d, err := dis.New("Bot " + c.Token)
	if err != nil {
		return nil, err
	}
	d.Identify.Intents =
		dis.IntentsGuilds |
			dis.IntentsGuildMessages
	return &Discord{
		masters: c.Masters,
		s:       d,
		rconC:   rcon,
	}, nil
}

func (d *Discord) Start() error {
	err := d.s.Open()
	if err != nil {
		return err
	}
	err = d.initCommand()
	if err != nil {
		return err
	}
	d.s.UpdateGameStatus(0, "正常に作動中...")
	return nil
}

func (d *Discord) Close() error {
	return d.s.Close()
}
