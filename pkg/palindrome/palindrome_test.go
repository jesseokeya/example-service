package palindrome

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPalindromeStrict(t *testing.T) {
	testCases := []struct {
		name string
		s    string
		want bool
	}{
		{
			"empty string",
			"",
			true,
		},
		{
			"single character",
			"a",
			true,
		},
		{
			"palindrome",
			"racecar",
			true,
		},
		{
			"palindrome with special character",
			"racecar!",
			false,
		},
		{
			"palindrome with upper case",
			"Racecar",
			false,
		},
		{
			"palindrome with whitespace",
			"a toyota",
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := IsPalindromeStrict(tc.s)
			require.Equal(t, tc.want, b)
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	testCases := []struct {
		name string
		s    string
		want bool
	}{
		{
			"empty string",
			"",
			true,
		},
		{
			"single character",
			"a",
			true,
		},
		{
			"palindrome",
			"racecar",
			true,
		},
		{
			"palindrome with special character",
			"racecar!",
			true,
		},
		{
			"palindrome with upper case",
			"Racecar",
			true,
		},
		{
			"palindrome with whitespace",
			"a toyota",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := IsPalindrome(tc.s)
			require.Equal(t, tc.want, b)
		})
	}
}
