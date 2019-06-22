package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	//args := os.Args[1:]

	ln, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening :", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Listening on ", CONN_HOST+":"+CONN_PORT)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting ", err.Error())
			continue
		}
		fmt.Println("Accepted connexion ", conn)
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	request := make([]byte, 1024)

	for {
		_, err := conn.Read(request)
		if err != nil {
			fmt.Println("Error reading : ", err.Error())
		}
		conn.Write(request)

		if request[0] == 'e' {
			conn.Write([]byte("exit\n"))
			break
		} else {
			// parse here
			conn.Write([]byte("Message recieved\n"))
		}
	}
}
