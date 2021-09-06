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

	*TLS Option:
	 	-- NOTE: for TLS to work the parameters to the LoadX509KeyPair function need to be replaced with the appropriate path of the key pair
	*Raw TCP Option:
	 	-- uncommment lines: 50
	 	-- comment lines: 36 -> 40 and 51

*/

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
)

func main() {

	// -- Server instance initialized
	s := newServer()

	// -- Loading KeyPair cert & key files

	cert, err := tls.LoadX509KeyPair("path-to-certificate-pem-file", "path-to-server's-key-pem-file")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// -- Creating listener
	PORT := ":" + arguments[1]
	// l, err := net.Listen("tcp4", PORT) // -- Raw TCP option
	l, err := tls.Listen("tcp4", PORT, config) // -- TLS option
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
