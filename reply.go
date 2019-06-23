package main

import (
	"fmt"
	"log"
)

func (session *Session) reply(message string) {
	session.Conn.Write([]byte(message))
	log.Printf("reply: '%s'", message)
}

func (session *Session) welcome() {
	message := fmt.Sprintf(":%s 001 %s :Welcome to the Internet Relay Network\r\n", CONN_HOST, session.Account.Nickname)
	session.reply(message)
}

func (session *Session) error401(name string) {
	message := fmt.Sprintf(":%s 401 %s :No such nick/channel\r\n", CONN_HOST, name)
	session.reply(message)
}

func (session *Session) error431() {
	message := fmt.Sprintf(":%s 431 :No nickname is given\r\n", CONN_HOST)
	session.reply(message)
}

func (session *Session) error432(nick string) {
	message := fmt.Sprintf(":%s 432 %s :Erroneus nickname\r\n", CONN_HOST, nick)
	session.reply(message)
}

func (session *Session) error433(nick string) {
	message := fmt.Sprintf(":%s 433 * %s :Nickname collision KILL\r\n", CONN_HOST, nick)
	session.reply(message)
}

func (session *Session) error434(nick string) {
	message := fmt.Sprintf(":%s 434 * %s :Nickname is already in use\r\n", CONN_HOST, nick)
	session.reply(message)
}

// Authorization

func (session *Session) error451() {
	message := fmt.Sprintf(":%s 451 :You have not registered\r\n", CONN_HOST)
	session.reply(message)
}

func (session *Session) error461(command string) {
	message := fmt.Sprintf(":%s 461 %s :Not enough parameters\r\n", CONN_HOST, command)
	session.reply(message)
}

func (session *Session) error462() {
	message := fmt.Sprintf(":%s 462 :You may not reregister\r\n", CONN_HOST)
	session.reply(message)
}

func (session *Session) error464() {
	message := fmt.Sprintf(":%s 464 :Password incorrect\r\n", CONN_HOST)
	session.reply(message)
}
