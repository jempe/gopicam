package validator

import (
	"errors"
	"regexp"
	"strconv"
)

const minLengthError string = "error:min_length_error"
const maxLengthError string = "error:max_length_error"
const alphaNumericAndDashesError string = "error:alphanumdash_error"

// Validate Username Length and make sure it contains only valid characters
func ValidateUsername(username string, minLength int, maxLength int) (valid bool, err error) {
	validMinLength, errMinLength := MinLength(username, minLength)
	if !validMinLength {
		err = errMinLength
		return
	}

	validMaxLength, errMaxLength := MaxLength(username, maxLength)
	if !validMaxLength {
		err = errMaxLength
		return
	}

	valid, errChars := AlphaNumericAndDashes(username)

	if !valid {
		err = errChars
		return
	}

	return
}

//Validate that input contains alphanumeric characters dashes or underscores only
func AlphaNumericAndDashes(value string) (valid bool, err error) {
	alphanumRegex := regexp.MustCompile(`^[a-z0-9-_]+$`)

	valid = alphanumRegex.MatchString(value)

	if !valid {
		err = errors.New(alphaNumericAndDashesError)
	}

	return
}

//Validate UUID format
func UUID(uuid string) (bool, error) {
	validUUID := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")

	if validUUID.MatchString(uuid) {
		return true, nil
	} else {
		err := errors.New("invalid_uuid")
		return false, err
	}
}

//Check if Length of string is at least a specific value
func MinLength(value string, min int) (bool, error) {
	if len(value) < min {
		return false, errors.New(minLengthError + "|" + strconv.Itoa(min))
	} else {
		return true, nil
	}
}

//Check if Length of string is at most a specific value
func MaxLength(value string, max int) (bool, error) {
	if len(value) > max {
		return false, errors.New(maxLengthError + "|" + strconv.Itoa(max))
	} else {
		return true, nil
	}
}

// Check if email is valid
func Email(email string) (bool, error) {
	validEmail := regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")

	if validEmail.MatchString(email) {
		return true, nil
	} else {
		err := errors.New("invalid_email")
		return false, err
	}
}
