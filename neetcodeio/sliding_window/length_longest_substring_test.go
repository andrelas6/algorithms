package slidingwindow

import (
	"math/rand"
	"strings"
	"testing"
)

// every test runs against both implementations in the file
var impls = []struct {
	name string
	fn   func(string) int
}{
	{"original", lengthOfLongestSubstring},
	{"optimized", lengthOfLongestSubstringOptimized},
}

func TestLengthOfLongestSubstring(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		// samples
		{"zxyzxyz", 3},
		{"xxxx", 1},
		{"abcabcbb", 3},
		{"bbbbb", 1},
		{"pwwkew", 3},

		// degenerate
		{"", 0},
		{"a", 1},
		{"ab", 2},
		{"aa", 1},

		// duplicate is the FIRST char of the window
		{"abba", 2},
		{"tmmzuxt", 5},
		{"abcdefga", 7},

		// printable ASCII, not just letters
		{" ", 1},
		{"  ", 1},
		{"a b a", 3},
		{"!@#!@#", 3},

		// answer sits at the very end
		{"aaabcdef", 6},
		// answer sits at the very start
		{"abcdefaaa", 6},
	}

	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			for _, c := range cases {
				if got := impl.fn(c.in); got != c.want {
					t.Errorf("%s(%q) = %d, want %d", impl.name, c.in, got, c.want)
				}
			}
		})
	}
}

// obviously-correct oracle
func brute(s string) int {
	best := 0
	for i := range s {
		seen := [128]bool{}
		for j := i; j < len(s); j++ {
			if seen[s[j]] {
				break
			}
			seen[s[j]] = true
			best = max(best, j-i+1)
		}
	}
	return best
}

func TestAgainstBruteForce(t *testing.T) {
	alphabets := []string{"ab", "abc", "abcde", "abcdefghij"}

	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			rng := rand.New(rand.NewSource(1))
			fails := 0
			for _, alpha := range alphabets {
				for trial := 0; trial < 2000; trial++ {
					n := rng.Intn(40)
					b := make([]byte, n)
					for i := range b {
						b[i] = alpha[rng.Intn(len(alpha))]
					}
					s := string(b)
					if got, want := impl.fn(s), brute(s); got != want {
						fails++
						if fails <= 3 {
							t.Errorf("%s(%q) = %d, want %d", impl.name, s, got, want)
						}
					}
				}
			}
			if fails > 3 {
				t.Errorf("%s: %d total mismatches out of 8000", impl.name, fails)
			}
		})
	}
}

// worst cases at the real constraint: 50,000 chars.
func benchInputs() map[string]string {
	var printable strings.Builder
	for c := byte(32); c <= 126; c++ {
		printable.WriteByte(c)
	}
	block := printable.String() // 95 distinct printable ASCII

	repeated := strings.Repeat(block, 50000/len(block)+1)[:50000]

	rng := rand.New(rand.NewSource(2))
	randomFull := make([]byte, 50000)
	for i := range randomFull {
		randomFull[i] = block[rng.Intn(len(block))]
	}

	return map[string]string{
		"all_same":        strings.Repeat("a", 50000),
		"two_chars":       strings.Repeat("ab", 25000),
		"cycle_95_alpha":  repeated,
		"random_95_alpha": string(randomFull),
		"all_distinct":    block + strings.Repeat("\x01", 0),
	}
}

func TestWorstCaseRuntime(t *testing.T) {
	for name, in := range benchInputs() {
		t.Run(name, func(t *testing.T) {
			for _, impl := range impls {
				t.Logf("%s: n=%d -> %d", impl.name, len(in), impl.fn(in))
			}
		})
	}
}

func BenchmarkLengthOfLongestSubstring(b *testing.B) {
	for _, impl := range impls {
		for name, in := range benchInputs() {
			b.Run(impl.name+"/"+name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					impl.fn(in)
				}
			})
		}
	}
}
