package main

import "regexp"

func isValidNickname(nick string) bool {
	r := regexp.MustCompile(`^[a-zA-Z0-9\-\[\]\\\x60\^\{\}]{1,8}$`)
	match := r.MatchString(nick)
	if match {
		return true
	}
	return false
}

func isValidUser(user string) bool {
	r := regexp.MustCompile(`[^\x20\x0\xd\xa]+`)
	match := r.MatchString(user)
	if match {
		return true
	}
	return false
}
