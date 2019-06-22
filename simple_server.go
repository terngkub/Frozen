package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

type Channel struct {
	Name string
	Topic string
	UserList []*Account
}

type Env struct {
	AccountList []Account
	UserList    map[string]*Account
	ConnList    map[string]net.Conn
}

type Session struct {
	Env     *Env
	Conn    net.Conn
	Account *Account
}

func main() {

	ln, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening :", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	env := Env{AccountList: []Account{}, UserList: make(map[string]*Account), ConnList: make(map[string]net.Conn)}

	fmt.Println("Listening on ", CONN_HOST+":"+CONN_PORT)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting ", err.Error())
			continue
		}
		fmt.Println("Accepted connexion ", conn)
		go runSession(&env, conn)
	}
}

func runSession(env *Env, conn net.Conn) {
	session := Session{Env: env, Conn: conn, Account: nil}
	defer session.Conn.Close()
	session.authorize()

	fmt.Println("UserList", session.Env.UserList)
	fmt.Println("SessionList", session.Env.ConnList)

	for {
		request, err := session.getRequest()
		if err != nil {
			break
		}
		session.handleRequest(request)
	}
}

func (session *Session) handleRequest(request string) {
}

func (session *Session) getRequest() (string, error) {
	request := make([]byte, 1024)
	len, err := session.Conn.Read(request)
	if err != nil {
		log.Println("Error reading: ", err)
		return "", err
	}
	requestStr := string(request[:len])
	fmt.Println("<" + requestStr + ">")
	return requestStr, nil
}

func privateMSG(session *Session) {
	request := session.getRequest()
	matches := doRegexpSubmatch("PRIVMSG (.*)\r\n", request)
	//if dst exists
	var dst string
	if len(matches) > 0 {
		for _, usr := range session.Env.AccountList {
			if matches[1] == usr.Nickname {
				dst = usr.Nickname
				fmt.Println(dst)
			}
		}
	}
	//grab message

	if len(request) > 1 {
		i := strings.Index(request[1:], ":")
		if i != 0 {
			msg := request[i+1:]
			//get dst's connexion
			dst_conn := session.Env.ConnList[dst]
			//send message from src to dst
			dst_conn.Write([]byte(msg))
		}
	}
}
