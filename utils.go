package main

import "regexp"

func doRegexpSubmatch(format, str string) []string {
	r := regexp.MustCompile(format)
	matches := r.FindStringSubmatch(str)
	return matches
}

func remove_user(s []*Account, i int) []*Account {
	if len(s) <= 1 {
		return nil
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func remove_chan(s []*Channel, i int) []*Channel {
	if len(s) <= 1 {
		return nil
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
