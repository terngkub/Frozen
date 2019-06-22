package main

import (
	"errors"
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
	session.Conn.Write([]byte(":localhost 001 " + user + " :Welcome"))

	return nil
}
