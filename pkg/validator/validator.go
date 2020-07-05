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

func UUID(uuid string) (bool, error) {
	validUUID := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")

	if validUUID.MatchString(uuid) {
		return true, nil
	} else {
		err := errors.New("invalid_uuid")
		return false, err
	}
}

func MinLength(value string, min int) (bool, error) {
	if len(value) < min {
		return false, errors.New("min_length_error")
	} else {
		return true, nil
	}
}

func MaxLength(value string, max int) (bool, error) {
	if len(value) > max {
		return false, errors.New("max_length_error")
	} else {
		return true, nil
	}
}

func Email(email string) (bool, error) {
	validEmail := regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")

	if validEmail.MatchString(email) {
		return true, nil
	} else {
		err := errors.New("invalid_email")
		return false, err
	}
}
