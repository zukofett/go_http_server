package main

import (
	"regexp"
)

var allowedCharsRX = regexp.MustCompile("^[-!#$%&'*+.^_`|~0-9A-Za-z]+$")

func isValidMethod(method string) bool {
    return allowedCharsRX.MatchString(method)
}
