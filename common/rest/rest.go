package rest

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"net/url"
	"time"
)

const contentType = "application/json"

const basePath = "/v1/api/"
const (
	getServerInfo    = "info"
	getPlayerList    = "players"
	getServerMetrics = "metrics"
	announceMessage  = "announce"
	kickPlayer       = "kick"
	banPlayer        = "ban"
	unbanPlayer      = "unban"
	saveWorld        = "save"
	shutdownServer   = "shutdown"
	forceStopServer  = "stop"
)

type Rest struct {
	c *resty.Client
}

func New(httpurl, password string) *Rest {
	u, _ := url.JoinPath(httpurl, basePath)
	return &Rest{
		c: resty.New().
			SetBaseURL(u).
			SetRetryCount(3).SetTimeout(10*time.Second).
			SetBasicAuth("admin", password),
	}
}

type GetServerInfoRsp struct {
	Version     string `json:"version"`
	Servername  string `json:"servername"`
	Description string `json:"description"`
}

func (r *Rest) GetServerInfo() (*GetServerInfoRsp, error) {
	rsp, err := r.c.R().Get(getServerInfo)
	if err != nil {
		return nil, err
	}
	var rspT GetServerInfoRsp
	err = parseRsp(rsp, &rspT)
	if err != nil {
		return nil, err
	}
	return &rspT, nil
}

type GetPlayerListRsp struct {
	Players []struct {
		Name      string  `json:"name"`
		Playerid  string  `json:"playerid"`
		Userid    string  `json:"userid"`
		Ip        string  `json:"ip"`
		Ping      float64 `json:"ping"`
		LocationX float64 `json:"location_x"`
		LocationY float64 `json:"location_y"`
		Level     int     `json:"level"`
	} `json:"players"`
}

func (r *Rest) GetPlayerList() (*GetPlayerListRsp, error) {
	rsp, err := r.c.R().Get(getPlayerList)
	if err != nil {
		return nil, err
	}
	var rspT GetPlayerListRsp
	err = parseRsp(rsp, &rspT)
	if err != nil {
		return nil, err
	}
	return &rspT, nil
}

type GetServerMetrics struct {
	Serverfps        int     `json:"serverfps"`
	Currentplayernum int     `json:"currentplayernum"`
	Serverframetime  float64 `json:"serverframetime"`
	Maxplayernum     int     `json:"maxplayernum"`
	Uptime           int     `json:"uptime"`
}

func (r *Rest) GetServerMetrics() (*GetServerMetrics, error) {
	rsp, err := r.c.R().Get(getServerMetrics)
	if err != nil {
		return nil, err
	}
	var rspT GetServerMetrics
	err = parseRsp(rsp, &rspT)
	if err != nil {
		return nil, err
	}
	return &rspT, nil
}

type announceMessageReq struct {
	Message string `json:"Message"`
}

func (r *Rest) AnnounceMessage(msg string) error {
	rsp, err := r.c.R().
		ForceContentType(contentType).
		SetBody(announceMessageReq{Message: msg}).
		Post(announceMessage)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New("request failed")
	}
	return nil
}

type kickPlayerReq struct {
	UserId  string `json:"userid"`
	Message string `json:"message"`
}

func (r *Rest) KickPlayer(playerId, msg string) error {
	rsp, err := r.c.R().SetBody(&kickPlayerReq{
		UserId:  playerId,
		Message: msg,
	}).
		ForceContentType(contentType).
		Post(kickPlayer)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New("request failed")
	}
	return nil
}

type banPlayerReq struct {
	UserId  string `json:"userid"`
	Message string `json:"message"`
}

func (r *Rest) BanPlayer(player, msg string) error {
	rsp, err := r.c.R().SetBody(&banPlayerReq{
		UserId:  player,
		Message: msg,
	}).
		ForceContentType(contentType).
		Post(banPlayer)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New("request failed")
	}
	return nil
}

type unBanPlayerReq struct {
	UserId string `json:"userid"`
}

func (r *Rest) UnbanPlayer(player string) error {
	rsp, err := r.c.R().SetBody(&unBanPlayerReq{
		UserId: player,
	}).
		ForceContentType(contentType).
		Post(unbanPlayer)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New("request failed")
	}
	return nil
}

func (r *Rest) SaveWorld() error {
	rsp, err := r.c.R().Get(saveWorld)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New("request failed")
	}
	return nil
}

type shutdownServerReq struct {
	WaitTime int    `json:"waittime"`
	Message  string `json:"message"`
}

func (r *Rest) ShutdownServer(wait int, msg string) error {
	rsp, err := r.c.R().SetBody(&shutdownServerReq{
		WaitTime: wait,
		Message:  msg,
	}).
		ForceContentType(contentType).
		Post(shutdownServer)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return errors.New("request failed")
	}
	return nil
}
