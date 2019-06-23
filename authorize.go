package main

import (
	"errors"
	"fmt"
)

func (account *Account) isComplete() bool {
	if account.Password != "" && account.Nickname != "" && account.User != "" {
		return true
	}
	return false
}

func (session *Session) authorize() error {
	newAccount := Account{}

	for !newAccount.isComplete() {

		request, err := session.getRequest()
		if err != nil {
			return err
		}

		newPass, err := session.cmdPASS(request)
		if err != nil {
			continue
		}
		if newPass != "" {
			newAccount.Password = newPass
		}

		matches := doRegexpSubmatch("NICK +(.+)\r\n", request)
		if len(matches) == 2 {
			_, duplicated := session.Env.NicknameMap[matches[1]]
			if duplicated {
				message := fmt.Sprintf(":%s 443 * %s :Nickname is already in use.\r\n", "127.0.0.1", matches[1])
				session.Conn.Write([]byte(message))
			} else {
				newAccount.Nickname = matches[1]
			}
		}

		matches = doRegexpSubmatch("USER +(.+) +.+ +.+ +:.+\r\n", request)
		if len(matches) == 2 {
			newAccount.User = matches[1]
		}
	}

	account, ok := session.Env.UserMap[newAccount.User]
	if ok {
		// check password
		if newAccount.Password != account.Password {
			return errors.New("wrong password")
		}
	} else {
		// create new user
		session.Env.AccountList = append(session.Env.AccountList, newAccount)
		session.Account = &session.Env.AccountList[len(session.Env.AccountList)-1]
		session.Env.UserMap[newAccount.User] = &session.Env.AccountList[len(session.Env.AccountList)-1]
		session.Env.NicknameMap[newAccount.Nickname] = &session.Env.AccountList[len(session.Env.AccountList)-1]
	}
	session.Env.ConnMap[newAccount.Nickname] = session.Conn
	message := fmt.Sprintf(":%s %s %s %s\r\n", "127.0.0.1", "001", newAccount.Nickname, ":Welcome to the Internet Relay Network")
	fmt.Println(message)
	session.Conn.Write([]byte(message))

	return nil
}

func (session *Session) cmdPASS(request string) (string, error) {
	if session.Account != nil {
		message := fmt.Sprintf(":%s 462 :You may not reregister", CONN_HOST)
		session.Conn.Write([]byte(message))
		return "", errors.New("462")
	}

	matches := doRegexpSubmatch("(?:.+ ){0,1}PASS +(.*)\r\n", request)
	if len(matches) != 2 {
		return "", nil
	}

	if matches[1] == "" {
		message := fmt.Sprintf(":%s 461 %s :Not enough parameters", CONN_HOST, "PASS")
		session.Conn.Write([]byte(message))
		return "", errors.New("461")
	}

	return matches[1], nil
}

func (session *Session) changeNickname(request string) {
	matches := doRegexpSubmatch("NICK +(.+)\r\n", request)
	if len(matches) != 2 {
		return
	}
	session.Env.NicknameMap[matches[1]] = session.Account
	session.Env.ConnMap[matches[1]] = session.Conn
	delete(session.Env.NicknameMap, session.Account.Nickname)
	delete(session.Env.ConnMap, session.Account.Nickname)
	session.Account.Nickname = matches[1]
}
