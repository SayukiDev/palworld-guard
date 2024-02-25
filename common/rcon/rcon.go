package rcon

import (
	"fmt"
	"strings"

	"github.com/gorcon/rcon"
)

// CmdName https://tech.palworldgame.com/server-commands
type CmdName string

const (
	Shutdown         CmdName = "Shutdown"         // /Shutdown {Seconds} {MessageText}
	DoExit           CmdName = "DoExit"           // /DoExit
	Broadcast        CmdName = "Broadcast"        // /Broadcast {MessageText}
	KickPlayer       CmdName = "KickPlayer"       // /KickPlayer {SteamID}
	BanPlayer        CmdName = "BanPlayer"        // /BanPlayer {SteamID}
	TeleportToPlayer CmdName = "TeleportToPlayer" // /TeleportToPlayer {SteamID}
	TeleportToMe     CmdName = "TeleportToMe"     // /TeleportToMe {SteamID}
	ShowPlayers      CmdName = "ShowPlayers"      // /ShowPlayers
	Info             CmdName = "Info"             // /Info
	Save             CmdName = "Save"             // /Save
)

type PalRcon struct {
	conn *rcon.Conn
}

func NewPalRcon(addr, password string) (*PalRcon, error) {
	conn, err := rcon.Dial(addr, password)
	if err != nil {
		return nil, err
	}
	return &PalRcon{
		conn: conn,
	}, nil
}

func (c *PalRcon) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("close conn error: %w", err)
	}
	return nil
}

func (c *PalRcon) Shutdown(seconds, message string) error {
	if _, err := c.exec(Shutdown, seconds, message); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) DoExit() error {
	if _, err := c.exec(DoExit); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) Broadcast(message string) error {
	if _, err := c.exec(Broadcast, message); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) KickPlayer(steamId string) error {
	if _, err := c.exec(KickPlayer, steamId); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) BanPlayer(steamId string) error {
	if _, err := c.exec(BanPlayer, steamId); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) TeleportToPlayer(steamId string) error {
	if _, err := c.exec(TeleportToPlayer, steamId); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) TeleportToMe(steamId string) error {
	if _, err := c.exec(TeleportToMe, steamId); err != nil {
		return err
	}
	return nil
}

type Player struct {
	Id   string
	Name string
}

func (c *PalRcon) ShowPlayers() ([]Player, error) {
	out, err := c.exec(ShowPlayers)
	if err != nil {
		return nil, err
	}
	raws := strings.Split(out, "\n")[1:]
	players := make([]Player, 0, len(raws))
	for _, raw := range raws {
		fields := strings.Split(raw, ",")
		if len(fields) < 2 {
			continue
		}
		players = append(players, Player{
			Id:   fields[1],
			Name: fields[0],
		})
	}
	return players, nil
}

func (c *PalRcon) Info() (string, error) {
	out, err := c.exec(Info)
	if err != nil {
		return "", err
	}
	return out, err
}

func (c *PalRcon) Save() error {
	if _, err := c.exec(Save); err != nil {
		return err
	}
	return nil
}

func (c *PalRcon) exec(cmd CmdName, args ...string) (string, error) {
	argStr := strings.Join(args, " ")
	cmdStr := string(cmd)
	if argStr != "" {
		cmdStr = fmt.Sprintf("%s %s", cmd, argStr)
	}
	out, err := c.conn.Execute(cmdStr)
	if err != nil {
		return "", fmt.Errorf("exec comnmand error: %w", err)
	}
	return out, nil
}
