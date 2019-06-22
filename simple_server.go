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

	fmt.Println("Listening on ", CONN_HOST+":"+CONN_PORT)
	defer ln.Close()
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
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading : ", err.Error())
	}
	conn.Write([]byte("Message recieved\n"))
	conn.Close()
}
