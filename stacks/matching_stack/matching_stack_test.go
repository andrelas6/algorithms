package matchingstack

import (
	"math/rand"
	"strings"
	"testing"
)

func TestAreBracketsProperlyMatched(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		// examples and samples from the statement
		{"example", "if (a[0] > b[1]) { doSomething(); }", true},
		{"no brackets at all", "int x = 42; // no brackets here", true},
		{"three pairs, separated", "() {} []", true},

		// degenerate
		{"empty string", "", true},
		{"only whitespace", "   ", true},

		// unclosed openers -> the classic missed check at the end
		{"single opener", "(", false},
		{"opener buried in text", "func main() {", false},
		{"all openers", "([{", false},

		// closer with nothing open
		{"single closer", ")", false},
		{"closer first", ")(", false},
		{"extra closer at the end", "()]", false},

		// wrong type: counts balance, nesting does not
		{"interleaved", "([)]", false},
		{"interleaved 2", "{[}]", false},
		{"mismatched pair", "(]", false},

		// properly nested
		{"nested three deep", "{[()]}", true},
		{"nested and sequential", "([]{})[]", true},
		{"deep nesting", "((((((((((()))))))))))", true},
		{"repeated pairs", "()()()()", true},

		// brackets mixed with other characters
		{"realistic code", "for (i := 0; i < len(x[0]); i++) { y[i] = z{a} }", true},
		{"text between brackets", "a(b[c]d)e", true},
		{"whitespace inside", "( [ ] )", true},
	}

	for _, c := range cases {
		if got := areBracketsProperlyMatched(c.in); got != c.want {
			t.Errorf("%s: areBracketsProperlyMatched(%q) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

// Independent oracle: strip everything that isn't a bracket, then repeatedly delete
// adjacent matching pairs. A valid string collapses to nothing.
func reduceOracle(s string) bool {
	var b strings.Builder
	for _, c := range s {
		if strings.ContainsRune("()[]{}", c) {
			b.WriteRune(c)
		}
	}
	cur := b.String()
	for {
		next := cur
		next = strings.ReplaceAll(next, "()", "")
		next = strings.ReplaceAll(next, "[]", "")
		next = strings.ReplaceAll(next, "{}", "")
		if next == cur {
			return cur == ""
		}
		cur = next
	}
}

func TestAgainstReduceOracle(t *testing.T) {
	r := rand.New(rand.NewSource(11))
	alphabet := []rune("()[]{}ab ;")

	for trial := 0; trial < 20000; trial++ {
		n := r.Intn(12)
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteRune(alphabet[r.Intn(len(alphabet))])
		}
		in := b.String()

		want := reduceOracle(in)
		if got := areBracketsProperlyMatched(in); got != want {
			t.Fatalf("areBracketsProperlyMatched(%q) = %v, want %v", in, got, want)
		}
	}
}
