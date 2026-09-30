package taskscheduler

import (
	"strings"
	"testing"
)

func TestLeastInterval(t *testing.T) {
	cases := []struct {
		name  string
		tasks string
		n     int
		want  int
	}{
		// samples from the statement
		{"sample 1", "AAABBB", 2, 8},
		{"sample 2", "ACABDB", 1, 6},
		{"sample 3", "AAABBB", 3, 10},

		// degenerate inputs
		{"empty", "", 2, 0},
		{"single task", "A", 5, 1},
		{"no cooldown", "AAABBB", 0, 6},
		{"no cooldown, one type", "AAAA", 0, 4},

		// one task type: every gap is pure idle
		{"one type, two copies", "AA", 5, 7},
		{"one type, many copies", "AAAA", 2, 10},

		// enough distinct tasks to fill every gap, so no idle at all
		{"all distinct", "ABCDEF", 2, 6},
		{"gaps filled exactly", "AABBCC", 2, 6},
		{"more types than slots", "AAABBBCCCDDE", 2, 12},

		// one dominant task, the rest fill part of its gaps
		{"dominant task", "AAAAAABCDEFG", 2, 16},

		// several tasks tied for the max count
		{"tie for max", "AAABBBCC", 2, 8},
		{"three-way tie", "AAABBBCCC", 2, 9},
		{"tie, big cooldown", "AAABBB", 50, 104},

		// greedy must pick the most frequent, not just any ready task
		{"pick most frequent", "AAAAABBC", 2, 13},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := leastInterval([]byte(c.tasks), c.n)
			if got != c.want {
				t.Errorf("leastInterval(%q, %d) = %d, want %d", c.tasks, c.n, got, c.want)
			}
		})
	}
}

func TestLeastIntervalMaxSize(t *testing.T) {
	// 10^4 copies of one task with n = 100: (10^4-1)*101 + 1
	tasks := []byte(strings.Repeat("A", 10000))
	if got, want := leastInterval(tasks, 100), 9999*101+1; got != want {
		t.Errorf("max size: got %d, want %d", got, want)
	}
}
