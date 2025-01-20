package main 

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
	//logger
	"bytes"
	"log"
)

const (
	PORT ="49999"
)

type Peer struct {
	address string
	conn net.Conn
}

//Direct messaging TCP 
func sendMessage(address string, message string) {
	conn, err := net.Dial("tcp", address)
	panicIfErrPresent(err)

	messageBytes := []byte(message)
	_, err = conn.Write(messageBytes)
	panicIfErrPresent(err)

	logger.debug(fmt.Sprintf("Message %s has been sent", message))
}

func startServer() {
	listener, err 
}
func sendMessage(address string, message string) {
	conn, err := net.Dial("tcp", address)
	panicIfErrPresent(err)
   
	messageBytes := []byte(message)
	_, err = conn.Write(messageBytes)
	panicIfErrPresent(err)
   
	logger.debug(fmt.Sprintf("Message %s has been sent", message))
   }
   
//create logger 
