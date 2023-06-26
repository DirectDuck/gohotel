package types

import (
	"regexp"
)

func IsEmailValid(e string) bool {
	// Sourced from https://stackoverflow.com/a/67686133
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
