package main

import (
	"fmt"
	"strings"
)

func (session *Session) privateMSG(request string) {
	src_nick := session.Account.Nickname
	src_user := session.Account.User
	matches := doRegexpSubmatch("PRIVMSG (.*) :(.*)", request)
	//if dst exists
	var dst_nick string
	if len(matches) > 0 {
		dst_nick = session.Env.NicknameMap[matches[1]].Nickname
	}
	//grab message
	if len(request) > 1 {
		i := strings.Index(request[1:], ":")
		if i != 0 {
			//:<nick>!<user>@<host> PRIVMSG dest :msg
			msg := fmt.Sprintf(":%s!%s@%s PRIVMSG %s :%s", src_nick,
				src_user,
				CONN_HOST,
				dst_nick,
				request[i+1:])
			//get dst's connexion
			dst_conn := session.Env.ConnMap[dst_nick]
			//send message from src to dst
			dst_conn.Write([]byte(msg))
		}
	}
}

func (session *Session) joinChan(request string) {
	src_user := session.Account
	found := false
	var req_keys []string
	sp_matches := strings.Split(request, " ")
	if len(sp_matches) >= 2 {
		req_chans := strings.Split(sp_matches[1], ",")
		if len(sp_matches) >= 3 {
			req_keys = strings.Split(sp_matches[2], ",")
		}
		//fmt.Println(req_chans)
		//fmt.Println(req_keys)
		// loop over req_chan
		for i, req_chan := range req_chans {
			// loop over channels
			for _, channel := range session.Env.ChannelList {
				// if req_chan exists
				if channel.Name == req_chan && is_not_banned(src_user, channel) {
					// check key
					if channel.Key != "" {
						if len(req_keys) > i && len(req_keys[i]) > 0 {
							if req_keys[i] == channel.Key {
								channel.UserList = append(channel.UserList, src_user)
							}
						}
					} else {
						channel.UserList = append(channel.UserList, src_user)
					}
					found = true
					break
				}
			}
			// if channel doesn't exists
			if found == false {
				if len(req_keys) > i && len(req_keys[i]) > 0 {
					session.createChannel(req_chans[i], "", req_keys[i])
				} else {
					session.createChannel(req_chans[i], "", "")
				}
			}
		}
	}
	//fmt.Println(session.Env.ChannelList)
}

func is_not_banned(user *Account, channel Channel) bool {
	for _, banned := range channel.BanList {
		if user.Nickname == banned.Nickname {
			return true
		}
	}
	return false
}

func (session *Session) createChannel(name string, topic string, key string) {
	new_chan := Channel{
		Name:      name,
		Topic:     topic,
		Key:       key,
		AdminList: []*Account{session.Account},
		UserList:  []*Account{session.Account},
		BanList:   []*Account{}}
	session.Env.ChannelList = append(session.Env.ChannelList, new_chan)
}
