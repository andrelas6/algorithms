package arrays_lists

import "testing"

// Spec under test, from the notes at the top of is_unique.go:
//   - only letters count towards uniqueness
//   - digits and symbols are "not valid" -> ignored entirely
//   - no case-folding is described, so 'a' and 'A' are distinct characters
func TestIsUnique(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		// degenerate inputs
		{"", true},
		{"a", true},

		// plain letters, no repeats
		{"abc", true},
		{"abcdefghijklmnopqrstuvwxyz", true},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", true},

		// plain letters, with repeats
		{"aa", false},
		{"abca", false},
		{"abcdefghijklmnopqrstuvwxyza", false},
		{"zyxz", false},

		// repeat is at the very end / very start
		{"abcc", false},
		{"aabc", false},

		// case sensitivity: upper and lower are different characters
		{"aA", true},
		{"Zz", true},
		{"aAbB", true},
		{"AaA", false},

		// digits and symbols are ignored, so repeats among them are fine
		{"11", true},
		{"%%%", true},
		{"$$^^", true},
		{"12313231%%%$$^^", true},
		{"1234567890", true},

		// ignored characters must not hide a real letter repeat
		{"a1a", false},
		{"a!!!a", false},
		{"a1b2c3a", false},
		{"%a%a%", false},

		// ignored characters interleaved with unique letters
		{"a1b2c3", true},
		{"!a@b#c$", true},
		{"1a2B3c4", true},

		// letters only unique once symbols/digits are stripped out
		{"h3ll0", false},  // two l's
		{"h3l0w", true},   // h, l, w
		{"a1b1c1d", true}, // 1 repeats but is ignored -> a,b,c,d unique
	}

	for _, c := range cases {
		if got := isUnique(c.in); got != c.want {
			t.Errorf("isUnique(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
