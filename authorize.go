package main

import (
	"errors"
	"fmt"
)

func (session *Session) authorize() error {
	var pass, nick, user string

	for pass == "" || nick == "" || user == "" {
		request, err := session.getRequest()
		if err != nil {
			return err
		}

		matches := doRegexpSubmatch("PASS (.*)\r\n", request)
		if len(matches) == 2 {
			pass = matches[1]
		}

		matches = doRegexpSubmatch("NICK (.*)\r\n", request)
		if len(matches) == 2 {
			nick = matches[1]
		}

		matches = doRegexpSubmatch("USER (.*) .* .* :.*\r\n", request)
		if len(matches) == 2 {
			user = matches[1]
		}
	}

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
	message := fmt.Sprintf(":%s %s %s %s\r\n", "127.0.0.1", "001", nick, ":Welcome to the Internet Relay Network")
	fmt.Println(message)
	session.Conn.Write([]byte(message))
	return nil
}
