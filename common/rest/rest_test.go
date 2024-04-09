package rest

import "testing"

var r = New("http://127.0.0.1:8212", "password")

func TestRest_GetPlayerList(t *testing.T) {
	t.Log(r.GetPlayerList())
}
