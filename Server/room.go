package main

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
