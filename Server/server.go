package main

/*
	- Server.go
	- This segment of the app is responsible for:
	  1- defining the structure of the server
	  2- outlining the ways in which client connections will be handeled.
	  3- outlining how each command from the client will be processed.

	  By: ahmed-alajbri

*/

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

//Server stores rooms [name,] & a channel where all comments will be sent through

type server struct {
	rooms map[string]*room
}

//Initializing new Server
func newServer() *server {
	return &server{
		rooms: make(map[string]*room),
	}
}

func (s *server) handleConnection(conn net.Conn) {
	// -- ACK New Client
	log.Printf("--[SERVER]: new client has joined: %s", conn.RemoteAddr().String())

	// -- Create client instance
	c := &client{
		conn: conn,
		name: "anonymous",
	}

	// -- List Commands :

	c.conn.Write([]byte("List of commands: \n" +
		"> /name -- set name \n" +
		"> /msg  -- send text message \n" +
		"> /join -- join existing or new chatroom \n" +
		"> /rooms -- list existing room \n" +
		"> /quit -- teriminate connection \n" +
		" example: /name Marco ,\n"))

	for {

		// -- Processing message from client -- waits for command from client
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// -- Message Parsing
		netData = strings.Trim(netData, "\r\n")
		args := strings.Split(netData, " ")
		cmd := strings.TrimSpace(args[0])

		// -- Command Validiation
		switch cmd {
		case "/name":
			s.setName(c, args)

		case "/join":
			s.joinRoom(c, args)

		case "/rooms":
			s.listRooms(c)

		case "/msg":
			s.msg(c, args)
		case "/quit":
			s.quit(c)
		default:
			c.err(fmt.Errorf("unknown command, please enter a valid command: %s", cmd))
		}

	}
	// conn.Close()
}

func (s *server) setName(c *client, args []string) {
	if len(args) < 2 {
		c.msg("name is required. usage: /name NAME")
		return
	}
	// -- Client name is set & acknowledgment sent to client
	c.name = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.name))
}

func (s *server) joinRoom(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	// -- Room existance validation
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	// -- Client added to room members
	r.members[c.conn.RemoteAddr()] = c

	// -- Remove Client from previous room
	s.quitCurrentRoom(c)
	c.room = r

	// -- Inform room that client has joined
	r.broadcast(c, fmt.Sprintf("%s joined the room", c.name))
	// -- Confirmation to the client
	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

// -- Method informs the client of the current avaialble rooms
func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	// -- List rooms to client
	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

// -- broadcast client message to the room
func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.name+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("hope to see you soon =(")
	c.conn.Close()
}

// -- Method responsible for removing client from the current room
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}
