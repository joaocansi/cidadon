package utils

import "strings"

func Mask(token string, percentage int) string {
	if token == "" || percentage <= 0 {
		return token
	}

	if percentage > 100 {
		percentage = 100
	}

	maskLen := len(token) * percentage / 100
	if maskLen == 0 {
		return token
	}

	return token[:len(token)-maskLen] + strings.Repeat("*", maskLen)
}
