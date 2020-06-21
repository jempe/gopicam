package validator

import (
	"errors"
	"regexp"
)

const usernameLessThanError string = "Username should be at least 6 characters long"
const usernameGreaterThanError string = "Username should up to 25 characters long"
const usernameRegexError string = "Username should contains alphanumeric characters, scores and underscores only"

// Validate Username Length and make sure it contains only valid characters
func ValidateUsername(username string) (valid bool, err error) {

	if len(username) < 6 {
		err = errors.New(usernameLessThanError)
		return
	}

	if len(username) > 25 {
		err = errors.New(usernameGreaterThanError)
		return
	}

	usernameRegex := regexp.MustCompile(`^[a-z0-9-_]*$`)

	valid = usernameRegex.MatchString(username)

	if !valid {
		err = errors.New(usernameRegexError)
	}

	return
}
