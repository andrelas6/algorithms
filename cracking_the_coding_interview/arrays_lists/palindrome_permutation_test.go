package arrays_lists

import "testing"

// Spec under test, straight from the statement:
//   - "you can ignore casing"      -> 'A' and 'a' are the SAME character
//   - "ignore non-letter characters" -> digits, spaces and symbols do not count
//   - permutation of a palindrome  -> at most one letter may have an odd count
//
// Empty (and letter-free) input is treated as true here: a string with nothing in
// it reads the same backwards, same as {"", true} in strings/two_pointers.
func TestPalindromePermutation(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		// the sample from the book
		{"Tact coa", true},
		{"taco cat", true},
		{"atco cta", true},

		// degenerate inputs
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},

		// odd-count budget: zero or one is fine, two is not
		{"aab", true},
		{"aabb", true},
		{"aabbc", true},
		{"aabbcc", true},
		{"aabbccd", true},
		{"abcabc", true},
		{"abc", false},
		{"aaabbb", false},
		{"aabbcd", false},
		{"hello", false},

		// casing is ignored -> these fold to the same letter
		{"Aa", true},
		{"AaB", true},
		{"Racecar", true},
		{"RaceCar", true},
		{"AaBbCc", true},
		{"AB", false},
		{"AbC", false},

		// non-letters are ignored entirely
		{"a1a", true},
		{"a b a", true},
		{"!!!", true},
		{"12321", true},
		{"1 2 3", true},
		{"a1b", false},
		{"ab!!", false},

		// classic palindromic phrases
		{"never odd or even", true},
		{"A man, a plan, a canal, Panama", true},
		{"Was it a car or a cat I saw", true},

		// near misses
		{"tact coat", false},
		{"taco cats", false},
	}

	for _, c := range cases {
		if got := palindromePermutation(c.in); got != c.want {
			t.Errorf("palindromePermutation(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
