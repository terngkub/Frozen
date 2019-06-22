package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type Account struct {
	Password string
	User     string
	Nickname string
}

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
		requestLen, err := conn.Read(request)
		if err != nil {
			fmt.Println("Error reading : ", err.Error())
		}
		fmt.Print(string(request))

		if requestLen == 0 {
			break
		} else {
			responseMessage(conn, request)
		}
	}
}

func responseMessage(conn net.Conn, request []byte) {
	requestStr := string(request)
	PASS := regexp.MustCompile("PASS (.*)\r\n")
	NICK := regexp.MustCompile("NICK (.*)\r\n")
	USER := regexp.MustCompile("USER (.*)\r\n")
	switch {
	case PASS.MatchString(requestStr):
		fmt.Println("match pass")
	case NICK.MatchString(requestStr):
		fmt.Println("match nick")
	case USER.MatchString(requestStr):
		fmt.Println("match user")
	}
}
