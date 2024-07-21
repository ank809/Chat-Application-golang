package helpers

import (
	"errors"
	"regexp"
)

var indianPhoneNumberPattern = `^[6-9]\d{9}$`

func IsValidIndianPhoneNumber(phoneNumber string) (bool, error) {
	re := regexp.MustCompile(indianPhoneNumberPattern)
	if !re.MatchString(phoneNumber) {
		return false, errors.New("invalid phone number")
	}
	return true, nil
}
