package arrays

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func TestIsAnagram(t *testing.T) {
	cases := []struct {
		name string
		s, t string
		want bool
	}{
		// samples
		{"racecar/carrace", "racecar", "carrace", true},
		{"jar/jam", "jar", "jam", false},
		{"anagram/nagaram", "anagram", "nagaram", true},
		{"rat/car", "rat", "car", false},

		// degenerate
		{"both empty", "", "", true},
		{"one empty", "a", "", false},
		{"other empty", "", "a", false},
		{"single same", "a", "a", true},
		{"single different", "a", "b", false},

		// length differs -> early exit
		{"prefix", "ab", "abc", false},

		// same letters, different counts
		{"aab vs abb", "aab", "abb", false},
		{"aabb vs abbb", "aabb", "abbb", false},

		// repeats
		{"all same letter", "aaaa", "aaaa", true},
		{"reversed", "abcdef", "fedcba", true},

		// same length, disjoint letters
		{"disjoint", "abc", "xyz", false},
	}

	for _, c := range cases {
		if got := isAnagram(c.s, c.t); got != c.want {
			t.Errorf("%s: isAnagram(%q, %q) = %v, want %v", c.name, c.s, c.t, got, c.want)
		}
	}
}

func TestIsAnagramIsSymmetric(t *testing.T) {
	pairs := [][2]string{{"listen", "silent"}, {"aab", "abb"}, {"", "a"}, {"abc", "cba"}}
	for _, p := range pairs {
		if isAnagram(p[0], p[1]) != isAnagram(p[1], p[0]) {
			t.Errorf("not symmetric for %q / %q", p[0], p[1])
		}
	}
}

func sortedLetters(s string) string {
	b := strings.Split(s, "")
	sort.Strings(b)
	return strings.Join(b, "")
}

func TestIsAnagramAgainstSorting(t *testing.T) {
	rng := rand.New(rand.NewSource(12))
	for trial := 0; trial < 3000; trial++ {
		mk := func() string {
			b := make([]byte, rng.Intn(8))
			for i := range b {
				b[i] = byte('a' + rng.Intn(4))
			}
			return string(b)
		}
		s, u := mk(), mk()
		want := sortedLetters(s) == sortedLetters(u)
		if got := isAnagram(s, u); got != want {
			t.Fatalf("isAnagram(%q, %q) = %v, want %v", s, u, got, want)
		}
	}
}
