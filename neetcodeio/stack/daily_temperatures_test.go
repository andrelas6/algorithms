package stack

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// dailyTemperaturesBruteForce is an independent oracle: for each day, scan
// forward to the first strictly warmer day. O(n^2), no stack, so a bug in the
// push/pop bookkeeping cannot hide behind the same mistake here.
func dailyTemperaturesBruteForce(temperatures []int) []int {
	result := make([]int, len(temperatures))
	for i := range temperatures {
		for j := i + 1; j < len(temperatures); j++ {
			if temperatures[j] > temperatures[i] {
				result[i] = j - i
				break
			}
		}
	}
	return result
}

func TestDailyTemperaturesExamples(t *testing.T) {
	cases := []struct {
		name  string
		temps []int
		want  []int
	}{
		// samples from the statement
		{"neetcode sample 1", []int{30, 38, 30, 36, 35, 40, 28}, []int{1, 4, 1, 2, 1, 0, 0}},
		{"neetcode sample 2", []int{22, 21, 20}, []int{0, 0, 0}},
		{"leetcode sample 1", []int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{"leetcode sample 2", []int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{"leetcode sample 3", []int{30, 60, 90}, []int{1, 1, 0}},

		// the trace in your own comments
		{"trace from the comments", []int{30, 39, 30, 36, 40}, []int{1, 3, 1, 1, 0}},

		// degenerate
		{"single day", []int{50}, []int{0}},
		{"two days, warmer", []int{50, 51}, []int{1, 0}},
		{"two days, colder", []int{51, 50}, []int{0, 0}},
	}

	for _, c := range cases {
		if got := dailyTemperatures(slices.Clone(c.temps)); !slices.Equal(got, c.want) {
			t.Errorf("%s: dailyTemperatures(%v) = %v, want %v", c.name, c.temps, got, c.want)
		}
	}
}

// The trap: the answer is the next STRICTLY warmer day. An equal temperature
// does not count, so a pop on >= instead of > gets every one of these wrong.
func TestDailyTemperaturesEqualIsNotWarmer(t *testing.T) {
	cases := []struct {
		name  string
		temps []int
		want  []int
	}{
		{"all equal", []int{70, 70, 70}, []int{0, 0, 0}},
		{"equal then warmer", []int{70, 70, 71}, []int{2, 1, 0}},
		{"equal days separated by colder ones", []int{70, 60, 70, 80}, []int{3, 1, 1, 0}},
		{"plateau inside a rise", []int{50, 60, 60, 60, 61}, []int{1, 3, 2, 1, 0}},
	}

	for _, c := range cases {
		if got := dailyTemperatures(slices.Clone(c.temps)); !slices.Equal(got, c.want) {
			t.Errorf("%s: dailyTemperatures(%v) = %v, want %v", c.name, c.temps, got, c.want)
		}
	}
}

// Shapes that make a single warm day answer many earlier days at once, which is
// the case the inner pop loop exists for: stop after one pop and these break.
func TestDailyTemperaturesOneDayAnswersMany(t *testing.T) {
	cases := []struct {
		name  string
		temps []int
		want  []int
	}{
		{"falling then a spike", []int{60, 50, 40, 30, 70}, []int{4, 3, 2, 1, 0}},
		{"spike clears only the colder days", []int{80, 60, 50, 70, 90}, []int{4, 2, 1, 1, 0}},
		{"two waves", []int{50, 40, 60, 45, 35, 65}, []int{2, 1, 3, 2, 1, 0}},
		{"strictly rising", []int{31, 32, 33, 34, 35}, []int{1, 1, 1, 1, 0}},
		{"strictly falling", []int{35, 34, 33, 32, 31}, []int{0, 0, 0, 0, 0}},
		{"zigzag", []int{40, 30, 50, 30, 60, 30, 70}, []int{2, 1, 2, 1, 2, 1, 0}},
		{"warmest day first", []int{100, 30, 40, 50}, []int{0, 1, 1, 0}},
	}

	for _, c := range cases {
		if got := dailyTemperatures(slices.Clone(c.temps)); !slices.Equal(got, c.want) {
			t.Errorf("%s: dailyTemperatures(%v) = %v, want %v", c.name, c.temps, got, c.want)
		}
	}
}

// Random inputs against the brute force. Narrow ranges force lots of equal
// temperatures; the full 30..100 range makes long rising and falling runs.
func TestDailyTemperaturesAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(17, 18))

	for range 2000 {
		n := 1 + rng.IntN(40)
		spread := []int{2, 5, 71}[rng.IntN(3)]

		temps := make([]int, n)
		for i := range temps {
			temps[i] = 30 + rng.IntN(spread)
		}

		got, want := dailyTemperatures(slices.Clone(temps)), dailyTemperaturesBruteForce(temps)
		if !slices.Equal(got, want) {
			t.Fatalf("dailyTemperatures(%v) = %v, brute force = %v", temps, got, want)
		}
	}
}

// Every answer must point at a warmer day, with nothing warmer in between, and
// a 0 must mean no warmer day exists at all. Checked straight from the
// definition rather than against another implementation.
func TestDailyTemperaturesAnswersMeetTheDefinition(t *testing.T) {
	temps := []int{73, 74, 75, 71, 69, 72, 76, 73, 73, 80, 30, 100, 99, 30}

	got := dailyTemperatures(slices.Clone(temps))
	if len(got) != len(temps) {
		t.Fatalf("got %d answers for %d days", len(got), len(temps))
	}

	for i, wait := range got {
		if wait == 0 {
			for j := i + 1; j < len(temps); j++ {
				if temps[j] > temps[i] {
					t.Errorf("day %d (%d): answer 0, but day %d (%d) is warmer", i, temps[i], j, temps[j])
					break
				}
			}
			continue
		}

		target := i + wait
		if target >= len(temps) || temps[target] <= temps[i] {
			t.Errorf("day %d (%d): answer %d does not land on a warmer day", i, temps[i], wait)
			continue
		}
		for j := i + 1; j < target; j++ {
			if temps[j] > temps[i] {
				t.Errorf("day %d (%d): answer %d, but day %d (%d) is warmer and sooner", i, temps[i], wait, j, temps[j])
				break
			}
		}
	}
}

// The input is only read.
func TestDailyTemperaturesDoesNotModifyInput(t *testing.T) {
	temps := []int{73, 74, 75, 71, 69, 72, 76, 73}
	original := slices.Clone(temps)

	dailyTemperatures(temps)

	if !slices.Equal(temps, original) {
		t.Errorf("dailyTemperatures modified its input: got %v, want %v", temps, original)
	}
}

// Constraints allow n up to 10^5. A falling run is the worst case for the stack:
// every day waits on it until the very end, then one day pops all of them. The
// whole point of the stack is that this is still O(n) overall.
func TestDailyTemperaturesAtMaxSize(t *testing.T) {
	const n = 100_000
	const budget = time.Second

	// never rising: a staircase from 99 down to 30 (equal days stay on the stack
	// too), then one final day warmer than everything
	fallingThenSpike := make([]int, n)
	for i := range n - 1 {
		fallingThenSpike[i] = 99 - (i*70/(n-1))%70
	}
	fallingThenSpike[n-1] = 100

	rising := make([]int, n)
	for i := range rising {
		rising[i] = 30 + i*70/n
	}

	allEqual := make([]int, n)
	for i := range allEqual {
		allEqual[i] = 65
	}

	rng := rand.New(rand.NewPCG(19, 20))
	random := make([]int, n)
	for i := range random {
		random[i] = 30 + rng.IntN(71)
	}

	cases := []struct {
		name  string
		temps []int
	}{
		{"falling then a spike", fallingThenSpike},
		{"rising", rising},
		{"all equal", allEqual},
		{"random", random},
	}

	for _, c := range cases {
		start := time.Now()
		got := dailyTemperatures(slices.Clone(c.temps))
		elapsed := time.Since(start)

		if len(got) != n {
			t.Errorf("%s: got %d answers, want %d", c.name, len(got), n)
			continue
		}
		if elapsed > budget {
			t.Errorf("%s (n=%d): took %v, over the %v budget", c.name, n, elapsed.Round(time.Millisecond), budget)
		}
		t.Logf("%s (n=%d): %v", c.name, n, elapsed.Round(time.Millisecond))
	}

	// spot-check the answers on the shapes where they are easy to know
	if got := dailyTemperatures(slices.Clone(allEqual)); slices.ContainsFunc(got, func(w int) bool { return w != 0 }) {
		t.Error("all equal at max size: every answer should be 0")
	}
	if got := dailyTemperatures(slices.Clone(fallingThenSpike)); got[0] != n-1 || got[n-1] != 0 {
		t.Errorf("falling then a spike: day 0 waits %d (want %d), last day waits %d (want 0)", got[0], n-1, got[n-1])
	}
}
