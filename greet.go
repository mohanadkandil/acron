package main

import "strings"

func Greet(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "world"
	}
	return "hello, " + name
}


