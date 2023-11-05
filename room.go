package main

import (
	"github.com/gobwas/ws"
	"net"
)

type Room struct {
	name    string
	members map[net.Addr]*Client
}

func (r *Room) broadcast(sender *Client, msg string) {
	for addr, member := range r.members {
		if addr != sender.conn.RemoteAddr() {
			member.msg(ws.OpText, msg)
		}
	}
}
