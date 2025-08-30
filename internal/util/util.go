package util

import "strings"

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Contains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
