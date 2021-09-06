package main

/*
	- room.go
	- This segment of the app defines the structue of the chat room and the broadcats method.

*/
import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

// -- Broadcast Method
func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		// -- Exclude the sender
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}
