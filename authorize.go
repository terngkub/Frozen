package main

import (
	"log"
	"strings"
)

func (session *Session) getRequests() ([]string, bool) {
	request := make([]byte, 512)
	len, err := session.Conn.Read(request)
	if err != nil {
		log.Println("error reading:", err)
		return []string{}, false
	}
	log.Printf("receive: '%s'", request)

	str := string(request[:len])

	if str[len-2] != '\r' || str[len-1] != '\n' {
		log.Println("error: request doesn't end with \\r\\n")
		return []string{}, true
	}
	trimmed := strings.Trim(str, "\r\n")
	splitted := strings.Split(trimmed, "\r\n")
	return splitted, true
}

func (session *Session) authorize() bool {
	account, ok := session.getAccountData()
	if !ok {
		return false
	}
	if oldAccount, ok := session.Env.UserMap[account.User]; ok {
		return session.login(account, oldAccount)
	}
	return session.register(account)
}

func (session *Session) getAccountData() (*Account, bool) {
	account := Account{}
	for !account.isComplete() {
		requests, ok := session.getRequests()
		if !ok {
			return nil, false
		}
		for _, request := range requests {
			switch {
			case strings.HasPrefix(request, "PASS"):
				if newPass := session.cmdPASS(request); newPass != "" {
					account.Password = newPass
				}
			case strings.HasPrefix(request, "NICK"):
				if newNick := session.cmdNICK(request); newNick != "" {
					account.Nickname = newNick
				}
			case strings.HasPrefix(request, "USER"):
				if newUser := session.cmdUSER(request); newUser != "" {
					account.User = newUser
					if account.Password != "" {
						if _, ok := session.Env.UserMap[account.User]; ok {
							return &account, true
						}
					}
				}
			default:
				session.error451()
				continue
			}
			if account.isComplete() {
				break
			}
		}
	}
	return &account, true
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

func (session *Session) login(newAccount *Account, oldAccount *Account) bool {
	log.Println("attemp to login", session.Conn)
	if newAccount.Password != oldAccount.Password {
		session.error464()
		return false
	}
	session.Account = oldAccount
	session.Env.ConnMap[oldAccount.Nickname] = session.Conn
	log.Println("login success", session.Conn)
	session.welcome()
	return true
}

func (session *Session) register(account *Account) bool {
	log.Println("attemp to register", session.Conn)
	if _, isDuplicated := session.Env.NicknameMap[account.Nickname]; isDuplicated {
		session.error434(account.Nickname)
		return false
	}
	session.Env.AccountList = append(session.Env.AccountList, *account)
	length := len(session.Env.AccountList)
	last := &session.Env.AccountList[length-1]
	session.Account = last
	session.Env.UserMap[account.User] = last
	session.Env.NicknameMap[account.Nickname] = last
	session.Env.ConnMap[account.Nickname] = session.Conn
	log.Println("register success", session.Conn)
	session.welcome()
	return true
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
