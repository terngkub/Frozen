package main

import (
	"errors"
	"strings"
)

func (session *Session) authorize() error {
	newAccount := Account{}

	for !newAccount.isComplete() {

		// TODO parse this in a better manner
		request, err := session.getRequest()
		if err != nil {
			return err
		}

		switch {
		case strings.HasPrefix(request, "PASS"):
			if newPass := session.cmdPASS(request); newPass != "" {
				newAccount.Password = newPass
			}
		case strings.HasPrefix(request, "NICK"):
			if newNick := session.cmdNICK(request); newNick != "" {
				newAccount.Nickname = newNick
			}
		case strings.HasPrefix(request, "USER"):
			if newUser := session.cmdUSER(request); newUser != "" {
				newAccount.User = newUser
			}
			// TODO handle unmatch case
		}
	}

	account, ok := session.Env.UserMap[newAccount.User]
	if ok {
		if newAccount.Password != account.Password {
			return errors.New("wrong password")
		}
	} else {
		session.register(newAccount)
	}
	session.Env.ConnMap[newAccount.Nickname] = session.Conn

	return nil
}

func (account *Account) isComplete() bool {
	if account.Password != "" && account.Nickname != "" && account.User != "" {
		return true
	}
	return false
}

func (session *Session) cmdPASS(request string) string {
	if session.Account != nil {
		session.error462()
		return ""
	}
	matches := doRegexpSubmatch("PASS +(.+)\r\n", request)
	if len(matches) != 2 {
		session.error461("PASS")
		return ""
	}
	return matches[1]
}

func (session *Session) cmdNICK(request string) string {
	matches := doRegexpSubmatch("NICK +(.+)\r\n", request)
	if len(matches) != 2 {
		session.error431()
		return ""
	}
	if !isValidNickname(matches[1]) {
		session.error432(matches[1])
		return ""
	}
	if _, isDuplicated := session.Env.NicknameMap[matches[1]]; isDuplicated {
		session.error434(matches[1])
		return ""
	}
	return matches[1]
}

func (session *Session) cmdUSER(request string) string {
	if session.Account != nil {
		session.error462()
		return ""
	}
	matches := doRegexpSubmatch("USER +(.+) +.+ +.+ +:.+\r\n", request)
	if len(matches) != 2 {
		session.error461("USER")
		return ""
	}
	return matches[1]
}

func (session *Session) register(account Account) {
	session.Env.AccountList = append(session.Env.AccountList, account)
	length := len(session.Env.AccountList)
	last := &session.Env.AccountList[length-1]
	session.Account = last
	session.Env.UserMap[account.User] = last
	session.Env.NicknameMap[account.Nickname] = last
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
