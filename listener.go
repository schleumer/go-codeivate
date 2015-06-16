package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "13666"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, fmt.Sprintf("%s:%s", CONN_HOST, CONN_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Somebody connected")
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Read loop
	for {
		// Make a buffer to hold incoming data.
		// Read the incoming connection into a string.
		status, err := bufio.NewReader(conn).ReadString('\n')

		if err == io.EOF {
			fmt.Println("Client disconnected")
			break
		} else if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		status = strings.Replace(status, "\r", "", -1)
		status = strings.Replace(status, "\n", "", -1)
		if len(status) < 1 {
			continue
		}

		fmt.Printf("Message received %s\n", status)
		conn.Write([]byte("ok\n"))
	}
}
