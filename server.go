package main

import (
	"errors"
	"fmt"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"strings"
)

type Server struct {
	rooms    map[string]*Room
	Commands chan Command
}

func newServer() *Server {
	return &Server{
		rooms:    make(map[string]*Room),
		Commands: make(chan Command),
	}
}

func (s *Server) run() {
	for cmd := range s.Commands {
		switch cmd.id {
		case CMD_NICK:
			s.nickname(cmd.client, cmd.args, cmd.opCode)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args, cmd.opCode)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args, cmd.opCode)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args, cmd.opCode)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args, cmd.opCode)
		}
	}
}
func (s *Server) NewClient(Conn net.Conn) *Client {
	log.Printf("new client has connected : %s", Conn.RemoteAddr().String())
	err := wsutil.WriteServerMessage(Conn, ws.OpText,
		[]byte("USAGE \n"+
			"commands are \n"+
			"/nick <name> - get a name, otherwise user will stay anonymous. \n"+
			"/join <name> - join a room, if room doesn't exist, the new room will be created. User can be only in one room at the same time. \n"+
			"/rooms - show list of available rooms to join. \n"+
			"/msg\t<msg> - broadcast message to everyone in a room. \n"+
			"/quit - disconnects from the chat server. \n"))
	if err != nil {
		// handle error
		log.Println("server->newClient", err)
	}
	c := &Client{
		conn:          Conn,
		nick:          "anonymous",
		room:          nil,
		clientComands: s.Commands,
	}
	return c
	//c.ReadInput()
}
func (s *Server) nickname(c *Client, args []string, op ws.OpCode) {
	if len(args) < 2 {
		c.msg(op, "nick is required. usage: /nick NAME")
		return
	}
	c.nick = args[1]
	c.msg(op, fmt.Sprintf("all right, I will call you %s", c.nick))
}
func (s *Server) join(c *Client, args []string, op ws.OpCode) {
	if len(args) < 2 {
		c.msg(op, "room name is required. usage: /join ROOM_NAME")
		return
	}
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &Room{
			name:    roomName,
			members: make(map[net.Addr]*Client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c
	s.quotCurrentRoom(c)
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(op, fmt.Sprintf("welcome to %s", r.name))
}
func (s *Server) listRooms(c *Client, args []string, op ws.OpCode) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	if rooms != nil {
		c.msg(op, fmt.Sprintf("avalibel rooms are : %s", strings.Join(rooms, ",")))
	} else {
		c.msg(op, fmt.Sprintf("No rooms have been created yet"))
	}

}
func (s *Server) msg(c *Client, args []string, op ws.OpCode) {
	if len(args) < 2 {
		c.msg(op, "message is required, usage: /msg MSG")
		return
	}
	if c.room == nil {
		c.err(op, errors.New("you must join the room first"))
		return
	}
	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}
func (s *Server) quit(c *Client, args []string, op ws.OpCode) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())
	if c.room == nil {
		s.quotCurrentRoom(c)

		c.msg(op, "sad to see you go :(")

		err := c.conn.Close()
		if err != nil {
			log.Printf("server->quit :%s", err)
		}
	} else {
		s.quotCurrentRoom(c)
		c.msg(op, fmt.Sprintf("you leave the Room"))
		c.room = nil

	}

}
func (s *Server) quotCurrentRoom(c *Client) {
	if c.room != nil {
		/*delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))*/

		oldroom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldroom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
		/*if len(s.rooms[c.room.name].members) == 0 {
			s.rooms[oldroom.name].members = nil
			s.rooms[oldroom.name].name = ""
		}*/
	}
}
