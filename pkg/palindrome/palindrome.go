// Package palindrome implements utilities to check if a string is a palindrome.
package palindrome

import (
	"regexp"
	"strings"
)

// IsPalindromeStrict returns true if s is a palindrome.
// An empty string is a palindrome;
// a single character is a palindrome;
// a string x y z is a palindrome, if y is a palindrome and x is a character equal to z;
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

// IsPalindrome converts s to lowercase and removes non-alphanumeric characters and whitespace from s before calling IsPalindromeStrict(s).
func IsPalindrome(s string) bool {
	r := regexp.MustCompile("[^a-zA-Z0-9]+")
	s = r.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	return IsPalindromeStrict(s)
}
