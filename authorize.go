package main

import (
	"errors"
	"log"
	"strings"
)

func (session *Session) getRequests() []string {
	request := make([]byte, 512)
	len, err := session.Conn.Read(request)
	if err != nil {
		log.Println("error reading: ", err)
		return []string{}
	}
	log.Printf("recieve: '%s'", request)

	str := string(request[:len])

	if str[len-2] != '\r' || str[len-1] != '\n' {
		log.Println("error: request doesn't end with \\r\\n")
		return []string{}
	}
	trimmed := strings.Trim(str, "\r\n")
	splitted := strings.Split(trimmed, "\r\n")
	return splitted
}

func (session *Session) authorize() error {
	newAccount := Account{}
	for !newAccount.isComplete() {
		requests := session.getRequests()
		for _, request := range requests {
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
			if newAccount.isComplete() {
				break
			}
		}
	}

	account, ok := session.Env.UserMap[newAccount.User]
	if ok {
		if newAccount.Password != account.Password {
			// TODO handle error464
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
	matches := doRegexpSubmatch("^PASS +(.+)$", request)
	if len(matches) != 2 {
		session.error461("PASS")
		return ""
	}
	return matches[1]
}

func (session *Session) cmdNICK(request string) string {
	matches := doRegexpSubmatch("^NICK +(.+)$", request)
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
	matches := doRegexpSubmatch("^USER +(.+) +.+ +.+ +:.+$", request)
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
