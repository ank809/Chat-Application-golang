package helpers

import "unicode"

func CheckPassword(password string) (bool, string) {
	if password == "" {
		return false, "Password cannot be empty"
	}
	if len(password) < 6 {
		return false, "Length of password should be greater than 6"
	}

	containsUpper := false
	containsLower := false
	containsDigits := false
	containsSpecialCharacters := false

	for _, ch := range password {
		if unicode.IsDigit(ch) {
			containsDigits = true
		} else if unicode.IsLower(ch) {
			containsLower = true
		} else if unicode.IsUpper(ch) {
			containsUpper = true
		} else {
			containsSpecialCharacters = true
		}
	}
	if containsDigits && containsLower && containsSpecialCharacters && containsUpper {
		return true, "Password is valid"
	} else {
		return false, "Password should contains uppercase, lowercase, digits and special characters"
	}
}
