package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Peer struct {
	address string
	conn    net.Conn
}

func main() {
	go startServer()

	scanner := bufio.NewScanner(os.Stdin)
	var currentPeer *Peer

	for scanner.Scan() {
		input := scanner.Text()
		if strings.HasPrefix(input, "/connect ") {
			addr := strings.TrimPrefix(input, "/connect")
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("Failed to connect to %s: %v\n", addr, err)
				continue
			}
			currentPeer = &Peer{address: addr, conn: conn}
			fmt.Printf("Connected to %s\n", addr)

			go handleConnection(conn) //start recieving msgs from peer
		} else if strings.HasPrefix(input, "/send ") && currentPeer != nil {
			message := strings.TrimPrefix(input, "/send ")
			_, err := currentPeer.conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			}

		} else if input == "/quit" {
			if currentPeer != nil {
				currentPeer.conn.Close()
			}
			os.Exit(0)
		} else {
			fmt.Println("Invalid command")
		}
	}
}

// create serve
func startServer() {
	listener, err := net.Listen("tcp", ":49999") //IPv4, IPv6
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1) //exit code for why program terminated
	}
	defer listener.Close() //ensure that listener is properly closed when the main() function exits
	fmt.Printf("Listening on port %s\n", ":49999")

	for {
		//connect to port
		conn, err := listener.Accept() //creates a conn object when a client connects to the server
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
			// log.Fatal(err)
		}

		fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())
		go handleConnection(conn)
		// go func (c net.Conn) {
		// 	io.Copy(c, c) //echo all incoming data
		// 	c.Close() // shut down conneciton
		// }(conn)
	}
}

func handleConnection(conn net.Conn) { //net.conn It is a type from Go's net package that represents a generic network connection
	defer conn.Close()
	reader := bufio.NewReader(conn) //'defer' ensures the connection is closed when the function returns // It runs even if the function encounters an error
	for {                           //infinite loop to continuously read messages
		message, err := reader.ReadString('\n') //ReadString reads until it encounters the delimiter '\n'
		if err != nil {
			fmt.Printf("Connection closed from %s\n", conn.RemoteAddr().String()) //RemoteAddr returns the remote network address
			return                                                                //automatically executes any deferred statements
		}
		fmt.Print("\nReceived from %s: %s", conn.RemoteAddr().String(), message)
	}
}

//An interface in Go defines a set of methods but doesn’t specify how they’re implemented
