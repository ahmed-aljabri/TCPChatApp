package main

/*
	- Client.go
	- This is a independant program that is responsible for initiating a TCP connection to the server
	  and runnung dedicated threads for both the Input and Output streams.

	  By: ahmed-alajbri

	  *Enabling TLS:
	  	-- Change the middle parameter to flag.String() on line 30 to the appropriate path of the cert
	*Raw TCP Option:
		-- comment lines: 34 -> 45 and 54
		-- uncomment lines: 55
		 NOTE: Make sure server and client operate on the same option (Raw-TCP/TLS)

*/
import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {

	certFile := flag.String("certfile", "/Path-to-CA-Certificate", "trusted CA certificate")
	flag.Parse()

	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", cert)
	}
	config := &tls.Config{RootCAs: certPool}

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := tls.Dial("tcp", CONNECT, config)
	// c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	// -- Recieve list of commands
	message, _ := bufio.NewReader(c).ReadString(',')
	fmt.Print("[SERVER]->: " + message)

	wg.Add(2)

	// -- Dedicated thread to the Input Stream
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

	// -- Dedicated thread to the Input Stream
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
