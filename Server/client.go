package main

import (
	"net"
)

type client struct {
	conn net.Conn
	name string
	room *room
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))

}
