package palindrome

import (
	"regexp"
	"strings"
)

// IsPalindromeStrict returns true if the provided string is a palindrome.
// The empty string is a palindrome;
// a string constituted only by a single character is a palindrome;
// a string c s d is a palindrome, if s is a palindrome and c is a character equal to d;
// nothing else is a palindrome.
func IsPalindromeStrict(s string) bool {
	l := len(s)
	for i := 0; i < l/2; i++ {
		if s[i] != s[l-i-1] {
			return false
		}
	}
	return true
}

// IsPalindrome returns true if the provided string is a palidrome.
// Only alphanumeric characters are considered;
// whitespace and case are ignored.
func IsPalindrome(s string) bool {
	r := regexp.MustCompile("[^a-zA-Z0-9]+")
	s = r.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	return IsPalindromeStrict(s)
}
