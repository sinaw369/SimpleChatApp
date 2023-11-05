package main

import "github.com/gobwas/ws"

const (
	CMD_NICK int = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type Command struct {
	id     int
	client *Client
	args   []string
	opCode ws.OpCode
}
