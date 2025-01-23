package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Message struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	go startServer() // Run the server in a goroutine
	startClient()    // Run the client in the main routine
}

// Start a server to listen for incoming messages
func startServer() {
	fmt.Println("What port are you listening on?")
	listener, err := net.Listen("tcp", ":8080") // Listen on port 8080
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Println("Someone has connected with you")
		go handleConnection(conn) // Handle each connection in a goroutine
	}
}

// Handle incoming connections and display received messages
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Friend has left the room: %v", err)
			return
		}
		log.Printf("Received: %s", message)
	}
}

func startClient() {
	log.Println("Starting client...")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Enter the peer's address (e.g., 127.0.0.1:8081):")
		scanner.Scan()
		address := scanner.Text()

		// Validate that the address includes a colon
		if !strings.Contains(address, ":") {
			fmt.Println("Invalid address format. Please include both hostname and port (e.g., 127.0.0.1:8081).")
			continue
		}

		conn, err := net.Dial("tcp", address) // Connect to the peer
		if err != nil {
			log.Printf("Error connecting to peer: %v", err)
			fmt.Println("Failed to connect. Please try again.")
			continue
		}
		defer conn.Close()

		fmt.Println("Connected! Type your messages (type 'exit' to quit):")
		writer := bufio.NewWriter(conn)
		go handleConnection(conn) // Handle incoming messages in a goroutine

		// Read user input and send messages
		for scanner.Scan() {
			text := scanner.Text()
			_, err := writer.WriteString(text)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
			if text == "exit" {
				log.Println("Exiting client...")
				return
			}
			sendMessage(conn, text)
		}
		break
	}
}

func sendMessage(conn net.Conn, text string) {
	msg := Message{
		Text:      text,
		Sender:    "User1",    // Replace with actual sender ID
		Receiver:  "User2",    // Replace with actual receiver ID
		Type:      "chat",     // Message type (e.g., "chat")
		Timestamp: time.Now(), // Timestamp of the message
	}

	// Serialize the message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	// Send the JSON message
	_, err = conn.Write(append(jsonData, '\n')) // Append newline for easier parsing
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
