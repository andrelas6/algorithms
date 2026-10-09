package characterreplacementlongestsubstring

import (
	"math/rand"
	"testing"
)

func TestCharacterReplacement(t *testing.T) {
	cases := []struct {
		name string
		s    string
		k    int
		want int
	}{
		// samples from the statement
		{"sample 1", "ABAB", 2, 4},
		{"sample 2", "AABABBA", 1, 4},

		// degenerate
		{"empty", "", 0, 0},
		{"single char", "A", 0, 1},
		{"single char, k=1", "A", 1, 1},

		// k = 0: plain longest run of one letter
		{"all same", "AAAA", 0, 4},
		{"all different", "ABCD", 0, 1},
		{"run at the end", "BAAA", 0, 3},
		{"run at the start", "AAAB", 0, 3},
		{"short run then longer run", "AABAAAB", 0, 3},

		// k covers everything
		{"k equals length", "ABCD", 4, 4},
		{"k fills the minority", "ABBB", 2, 4},

		// best window is not where the early best was
		{"best window later in string", "ABCDEAAAAA", 2, 7},
		{"replace into the longer run", "AABBBCCCC", 1, 5},
		{"three letters, k=2", "ABCABCAAA", 2, 6},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := characterReplacement(c.s, c.k); got != c.want {
				t.Errorf("characterReplacement(%q, %d) = %d, want %d", c.s, c.k, got, c.want)
			}
		})
	}
}

// Compares against an O(n^2) brute force on random small strings.
func TestCharacterReplacementRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for range 2000 {
		n := rng.Intn(15)
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('A' + rng.Intn(3))
		}
		s := string(b)
		k := rng.Intn(n + 1)

		if got, want := characterReplacement(s, k), bruteForce(s, k); got != want {
			t.Fatalf("characterReplacement(%q, %d) = %d, want %d", s, k, got, want)
		}
	}
}

func bruteForce(s string, k int) int {
	best := 0
	for i := 0; i < len(s); i++ {
		var count [26]int
		top := 0
		for j := i; j < len(s); j++ {
			count[s[j]-'A']++
			top = max(top, count[s[j]-'A'])
			if j-i+1-top <= k {
				best = max(best, j-i+1)
			}
		}
	}
	return best
}
