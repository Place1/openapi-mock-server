package utils

import "strings"

// ContainsUpperChars returns true if the input
// string contains any upper-case letters.
func ContainsUpperChars(input string) bool {
	return strings.ToLower(input) != input
}
