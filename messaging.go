package main

import (
	"fmt"
	"strings"
)

func (session *Session) privateMSG(request string) {
	matches := doRegexpSubmatch("^PRIVMSG +(.+?) +:(.+)$", request)
	if len(matches) != 3 {
		session.error461("PRIVMSG")
		return
	}

	msg := fmt.Sprintf(":%s!%s@%s PRIVMSG %s :%s\r\n",
		session.Account.Nickname,
		session.Account.User,
		CONN_HOST,
		matches[1],
		matches[2])

	if matches[1][0] == '#' || matches[1][0] == '&' {
		if channel, ok := session.Env.ChannelMap[matches[1]]; ok {
			for _, account := range channel.UserList {
				if account.Nickname != session.Account.Nickname {
					dstConn := session.Env.ConnMap[account.Nickname]
					dstConn.Write([]byte(msg))
				}
			}
		} else {
			session.error401(matches[1])
		}
	} else {
		if account, ok := session.Env.NicknameMap[matches[1]]; ok {
			dstConn := session.Env.ConnMap[account.Nickname]
			dstConn.Write([]byte(msg))
		} else {
			session.error401(matches[1])
		}
	}
}

func (session *Session) joinChan(request string) {
	src_user := session.Account
	var req_keys []string
	sp_matches := strings.Split(request, " ")
	// grab channels and keys from request
	if len(sp_matches) >= 2 {
		req_chans := strings.Split(sp_matches[1], ",")
		if len(sp_matches) >= 3 {
			req_keys = strings.Split(sp_matches[2], ",")
		}
		for i, req_chan := range req_chans {
			channel, ok := session.Env.ChannelMap[req_chan]
			if ok == true {
				if is_banned(src_user, *channel) == false {
					// check key
					session.checkChan(channel, req_keys, i)
				}
			} else if len(req_keys) > i && len(req_keys[i]) > 0 {
				session.createChannel(req_chans[i], "", req_keys[i])
			} else {
				session.createChannel(req_chans[i], "", "")
			}
		}
	} else {
		session.error461(request)
	}
}

func (session *Session) checkChan(channel *Channel, req_keys []string, idx int) {
	if channel.Key != "" {
		if len(req_keys) > idx && len(req_keys[idx]) > 0 {
			if req_keys[idx] == channel.Key {
				session.append_user(channel)
			}
		}
	} else {
		session.append_user(channel)
	}
	//TODO message if banned

}

func is_banned(user *Account, channel Channel) bool {
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
		UserList:  []*Account{},
		BanList:   []*Account{},
		UserMap:   make(map[string]*Account)}
	session.Env.ChannelList = append(session.Env.ChannelList, &new_chan)
	session.Env.ChannelMap[name] = &new_chan
	session.append_user(&new_chan)
}

func (session *Session) append_user(channel *Channel) {
	// add user to channel
	src_user := session.Account
	channel.UserList = append(channel.UserList, src_user)
	channel.UserMap[src_user.Nickname] = src_user
	// alert other users
	for _, user := range channel.UserList {
		alert_message := fmt.Sprintf("%s!%s@%s %s %s\r\n",
			src_user.Nickname,
			src_user.User,
			CONN_HOST,
			"JOIN",
			channel.Name)
		session.Env.ConnMap[user.Nickname].Write([]byte(alert_message))
	}
	//send responses
	topic_message := fmt.Sprintf(":%s!%s@%s %s %s %s :%s\r\n",
		src_user.Nickname,
		src_user.User,
		CONN_HOST,
		"332",
		src_user.Nickname,
		channel.Name,
		channel.Topic)
	session.Conn.Write([]byte(topic_message))
	var users_list string
	for _, user := range channel.UserList {
		users_list += user.Nickname + " "
	}
	names_message := fmt.Sprintf(":%s!%s@%s %s %s = %s :%s\r\n",
		src_user.Nickname,
		src_user.User,
		CONN_HOST,
		"353",
		src_user.Nickname,
		channel.Name,
		users_list)
	session.Conn.Write([]byte(names_message))
	end_names_message := fmt.Sprintf(":%s!%s@%s %s %s %s :%s\r\n",
		src_user.Nickname,
		src_user.User,
		CONN_HOST,
		"366",
		src_user.Nickname,
		channel.Name,
		"End of NAMES list")
	session.Conn.Write([]byte(end_names_message))
}

func (session *Session) leaveChan(request string) {
	if len(request) > 5 {
		src := session.Account
		matches := doRegexpSubmatch("PART (.*) :(.*)", request)
		if len(matches) > 0 {
			// if channel/user exists
			channel, ok1 := session.Env.ChannelMap[matches[1]]
			if ok1 == true {
				_, ok2 := channel.UserMap[src.Nickname]
				if ok2 == true {
					session.sendPart(request, src, channel)
				}
			}
		}
	}
}
func (session *Session) sendPart(request string, src *Account, channel *Channel) {
	// send PART messages
	i := strings.Index(request[1:], ":")
	msg := fmt.Sprintf(":%s!%s@%s PART %s :%s\r\n",
		src.Nickname,
		src.User,
		CONN_HOST,
		channel.Name,
		request[i+2:])
	for _, user := range channel.UserList {
		if user.Nickname != src.Nickname {
			dst_conn := session.Env.ConnMap[user.Nickname]
			dst_conn.Write([]byte(msg))
		}
	}
	// leave chann
	for i, user := range channel.UserList {
		if user.Nickname == src.Nickname {
			channel.UserList = remove_user(channel.UserList, i)
			delete(channel.UserMap, user.Nickname)
		}
	}
	// remove empty channel
	if channel.UserList == nil {
		for i, channel := range session.Env.ChannelList {
			if channel.Name == channel.Name {
				session.Env.ChannelList = remove_chan(session.Env.ChannelList, i)
				delete(session.Env.ChannelMap, channel.Name)
			}
		}
	}
}
