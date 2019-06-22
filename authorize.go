package main

import (
	"errors"
	"fmt"
)

func (session *Session) authorize() error {
	request := session.getRequest()
	matches := doRegexpSubmatch("PASS (.*)\r\nNICK (.*)\r\nUSER (.*) .* .* :.*\r\n", request)
	if len(matches) != 4 {
		return errors.New("wrong request format")
	}
	pass := matches[1]
	nick := matches[2]
	user := matches[3]

	account, ok := session.Env.UserList[user]
	if ok {
		// check password
		if pass != account.Password {
			return errors.New("wrong password")
		}
	} else {
		// create new user
		newAccount := Account{Password: pass, Nickname: nick, User: user}
		session.Env.AccountList = append(session.Env.AccountList, newAccount)
		session.Account = &session.Env.AccountList[len(session.Env.AccountList)-1]
		session.Env.UserList[user] = &session.Env.AccountList[len(session.Env.AccountList)-1]
	}
	session.Env.ConnList[nick] = session.Conn
	fmt.Println("correct")
	message := fmt.Sprintf(":%s 001 %s :Welcome to the Internet Relay Network %s!%s@%s", "127.0.0.1", nick, nick, user, "127.0.0.1")
	fmt.Println(message)
	session.Conn.Write([]byte(message))

	return nil
}
