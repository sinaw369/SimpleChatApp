package main

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn          net.Conn
	nick          string
	room          *Room
	clientComands chan<- Command
}

func (c *Client) ReadInput() {
	for {
		msg, op, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			// handle error
			err := c.conn.Close()
			if err != nil {
				//log.Println("closeConnection err->readInput :", err)
				break
			}
			//log.Println("client->readinput :", err)

			//panic(err)
		}
		stMsg := string(msg)
		stMsg = strings.Trim(stMsg, "\r\n")
		args := strings.Split(stMsg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.clientComands <- Command{
				id:     CMD_NICK,
				client: c,
				args:   args,
				opCode: op,
			}
		case "/join":
			c.clientComands <- Command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
				opCode: op,
			}
		case "/rooms":
			c.clientComands <- Command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
				opCode: op,
			}
		case "/msg":
			c.clientComands <- Command{
				id:     CMD_MSG,
				client: c,
				args:   args,
				opCode: op,
			}
		case "/quit":
			c.clientComands <- Command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
				opCode: op,
			}
		default:
			c.err(op, fmt.Errorf("unknown command:%s", cmd))

		}
	}
}
func (c *Client) err(op ws.OpCode, err error) {
	connErr := []byte("Error:" + err.Error() + "\n")
	err = wsutil.WriteServerMessage(c.conn, op, connErr)
	if err != nil {
		// handle error
		log.Printf("client->err :%s", err)
	}
}
func (c *Client) msg(op ws.OpCode, msg string) {
	err := wsutil.WriteServerMessage(c.conn, op, []byte("> "+msg+"\n"))
	if err != nil {
		// handle error
		log.Printf("client->msg :%s", err)
	}
}
