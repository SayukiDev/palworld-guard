package config

import "testing"

var c *Config

func init() {
	c = New("../example/config.json5")
}

func TestConfig_Save(t *testing.T) {
	c.Save()
}

func TestConfig_Load(t *testing.T) {
	c.Load()
	t.Log(c)
	t.Log(c.Rcon)
	t.Log(c.Process)
	t.Log(c.Discord)
}
