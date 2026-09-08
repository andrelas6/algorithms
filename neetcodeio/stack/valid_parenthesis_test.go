package stack

import "testing"

func TestIsValid(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		// samples
		{"single pair", "[]", true},
		{"three kinds in a row", "([{}])", true},
		{"crossed", "[(])", false},
		{"simple pair", "()", true},
		{"all three flat", "()[]{}", true},
		{"mismatched pair", "(]", false},

		// degenerate
		{"empty", "", true},
		{"single open", "(", false},
		{"single close", ")", false},

		// closing with nothing on the stack
		{"close first", ")(", false},
		{"extra close at end", "()]", false},

		// unclosed openers left over
		{"unclosed", "([", false},
		{"unclosed nested", "([]", false},

		// deep nesting
		{"deep", "((((()))))", true},
		{"deep mismatch", "((((())))", false},

		// wrong closer for the top of the stack
		{"wrong closer", "{[}]", false},
		{"right kinds wrong order", "[{]}", false},
	}

	for _, c := range cases {
		if got := isValid(c.in); got != c.want {
			t.Errorf("%s: isValid(%q) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

func TestIsOpenBracketAndCo(t *testing.T) {
	for _, c := range []rune{'(', '[', '{'} {
		if !isOpenBracketAndCo(c) {
			t.Errorf("isOpenBracketAndCo(%q) = false, want true", c)
		}
	}
	for _, c := range []rune{')', ']', '}', 'a', '0', ' '} {
		if isOpenBracketAndCo(c) {
			t.Errorf("isOpenBracketAndCo(%q) = true, want false", c)
		}
	}
}
