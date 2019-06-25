package main

import "regexp"

func doRegexpSubmatch(format, str string) []string {
	r := regexp.MustCompile(format)
	matches := r.FindStringSubmatch(str)
	return matches
}

func remove(s []*Account, i int) []*Account {
	if len(s) <= 1 {
		return nil
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
