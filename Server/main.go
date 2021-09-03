package main

/*
	- main.go
	- This segment of the app is responsible for:
	  1- Creating an instance of the server.
	  2- Running a listenr on the port provided by the server initiator.
	  3- Listening for client connections
	  4- Creating dedicated threads for each client connection which is then Passed
	  	 to the server [handleConnection()] to start processing the clients.

	  By: ahmed-alajbri

*/

import (
	"fmt"
	"net"
	"os"
)

func main() {

	// -- Server instance initialized
	s := newServer()

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// -- Creating listener
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// -- Listening for client connections
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// - Dedicated client threads created and passed to the server for processing
		go s.handleConnection(c)

	}
}
