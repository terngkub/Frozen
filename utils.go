package main

import "regexp"

func doRegexpSubmatch(format, str string) []string {
	r := regexp.MustCompile(format)
	matches := r.FindStringSubmatch(str)
	return matches
}
