package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	//defer c.Close()

	// -- RECIEVE ACK OF CONNECTION & LIST OF COMMAND
	message, _ := bufio.NewReader(c).ReadString(',')
	fmt.Print("[SERVER]->: " + message)

	wg.Add(2)
	// -- Incoming Stream Thread
	go func(c net.Conn) {

		for {
			message, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				log.Fatalf("unable to read Input Stream: %s", err.Error())
			}
			fmt.Print("[SERVER]->: " + message)
		}

		wg.Done()
	}(c)

	// -- Outgowing Stream Thread
	go func(c net.Conn) {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println(">> ")
			text, _ := reader.ReadString('\n')
			fmt.Fprintf(c, text+"\n")

			if text == "/quit" {
				fmt.Println("TCP client exiting...")

				return
			}
		}

		wg.Done()
	}(c)

	wg.Wait()

	defer c.Close()
}
