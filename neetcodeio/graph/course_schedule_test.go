package graph

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// canFinishByIndegree is an independent oracle: Kahn's algorithm. Repeatedly take
// a course with no remaining prerequisites; if every course gets taken, there
// is no cycle. Iterative and count-based, no recursion and no visiting set, so a
// bug in the DFS cycle check cannot hide behind the same mistake here.
func canFinishByIndegree(numCourses int, prerequisites [][]int) bool {
	unlocks := make([][]int, numCourses) // prereq -> courses it unlocks
	remaining := make([]int, numCourses) // course -> prerequisites not yet taken

	for _, pair := range prerequisites {
		course, prereq := pair[0], pair[1]
		unlocks[prereq] = append(unlocks[prereq], course)
		remaining[course]++
	}

	ready := []int{}
	for course, n := range remaining {
		if n == 0 {
			ready = append(ready, course)
		}
	}

	taken := 0
	for len(ready) > 0 {
		course := ready[len(ready)-1]
		ready = ready[:len(ready)-1]
		taken++

		for _, next := range unlocks[course] {
			remaining[next]--
			if remaining[next] == 0 {
				ready = append(ready, next)
			}
		}
	}

	return taken == numCourses
}

func TestCanFinishExamples(t *testing.T) {
	cases := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		want          bool
	}{
		// samples from the statement
		{"neetcode sample 1", 2, [][]int{{0, 1}}, true},
		{"neetcode sample 2", 2, [][]int{{0, 1}, {1, 0}}, false},
		{"leetcode sample 1", 2, [][]int{{1, 0}}, true},
		{"leetcode sample 2", 2, [][]int{{1, 0}, {0, 1}}, false},

		// degenerate
		{"single course", 1, [][]int{}, true},
		{"no prerequisites at all", 5, [][]int{}, true},
		{"nil prerequisites", 3, nil, true},
		{"course is its own prerequisite", 1, [][]int{{0, 0}}, false},
		{"self loop among others", 3, [][]int{{1, 0}, {2, 2}}, false},
	}

	for _, c := range cases {
		if got := canFinish(c.numCourses, c.prerequisites); got != c.want {
			t.Errorf("%s: canFinish(%d, %v) = %v, want %v", c.name, c.numCourses, c.prerequisites, got, c.want)
		}
	}
}

// Cycles of different lengths and in different places. The cycle does not have
// to involve course 0, and it does not have to be reachable from the first
// course the loop starts at.
func TestCanFinishCycles(t *testing.T) {
	cases := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
	}{
		{"three-cycle", 3, [][]int{{0, 1}, {1, 2}, {2, 0}}},
		{"cycle away from course 0", 4, [][]int{{1, 0}, {2, 3}, {3, 2}}},
		{"cycle only reachable from the last course", 5, [][]int{{4, 3}, {3, 2}, {2, 3}}},
		{"long chain closed at the end", 6, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 0}}},
		{"cycle hanging off a valid chain", 5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 2}}},
		{"two separate components, one cyclic", 6, [][]int{{1, 0}, {2, 1}, {4, 3}, {5, 4}, {3, 5}}},
	}

	for _, c := range cases {
		if canFinish(c.numCourses, c.prerequisites) {
			t.Errorf("%s: canFinish(%d, %v) = true, want false (there is a cycle)", c.name, c.numCourses, c.prerequisites)
		}
	}
}

// The trap: two courses sharing a prerequisite is NOT a cycle. A cycle check
// that marks courses "seen" and never unmarks them reaches the shared course a
// second time and wrongly reports a cycle on every one of these.
func TestCanFinishSharedPrerequisitesAreNotCycles(t *testing.T) {
	cases := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
	}{
		// 0 needs 1 and 2, both of which need 3
		{"diamond", 4, [][]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}}},
		// 0 needs 1, 2 and 3 directly, and 1 and 2 also need 3
		{"shared prerequisite at two depths", 4, [][]int{{0, 1}, {0, 2}, {0, 3}, {1, 3}, {2, 3}}},
		// every course needs course 0
		{"one prerequisite for everything", 5, [][]int{{1, 0}, {2, 0}, {3, 0}, {4, 0}}},
		// one course needs everything else
		{"one course needs all the others", 5, [][]int{{0, 1}, {0, 2}, {0, 3}, {0, 4}}},
		// two diamonds stacked
		{"stacked diamonds", 7, [][]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}, {3, 5}, {4, 6}, {5, 6}}},
		// same prerequisite reached from separate starting courses
		{"two starts, one shared chain", 5, [][]int{{0, 2}, {1, 2}, {2, 3}, {3, 4}}},
	}

	for _, c := range cases {
		if !canFinish(c.numCourses, c.prerequisites) {
			t.Errorf("%s: canFinish(%d, %v) = false, want true (shared prerequisites, no cycle)", c.name, c.numCourses, c.prerequisites)
		}
	}
}

// Valid orderings of different shapes, including prerequisites listed out of
// order and courses that nothing refers to.
func TestCanFinishValidSchedules(t *testing.T) {
	cases := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
	}{
		{"chain", 5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}},
		{"chain listed backwards", 5, [][]int{{3, 4}, {2, 3}, {1, 2}, {0, 1}}},
		{"chain the other way", 5, [][]int{{4, 3}, {3, 2}, {2, 1}, {1, 0}}},
		{"isolated courses around a chain", 6, [][]int{{2, 3}, {3, 4}}},
		{"two independent chains", 6, [][]int{{0, 1}, {1, 2}, {3, 4}, {4, 5}}},
		{"both directions between different pairs", 4, [][]int{{0, 1}, {3, 2}}},
	}

	for _, c := range cases {
		if !canFinish(c.numCourses, c.prerequisites) {
			t.Errorf("%s: canFinish(%d, %v) = false, want true", c.name, c.numCourses, c.prerequisites)
		}
	}
}

// Random graphs against Kahn's algorithm. Sparse graphs are mostly acyclic,
// dense ones mostly cyclic, so both answers get plenty of coverage. Fixed seed
// keeps it repeatable.
func TestCanFinishAgainstIndegree(t *testing.T) {
	rng := rand.New(rand.NewPCG(21, 22))

	trueCount, falseCount := 0, 0
	for range 3000 {
		numCourses := 1 + rng.IntN(10)
		numEdges := rng.IntN(2 * numCourses)

		seen := map[[2]int]bool{}
		var prerequisites [][]int
		for range numEdges {
			pair := [2]int{rng.IntN(numCourses), rng.IntN(numCourses)}
			if seen[pair] {
				continue // pairs are unique per the constraints
			}
			seen[pair] = true
			prerequisites = append(prerequisites, []int{pair[0], pair[1]})
		}

		got := canFinish(numCourses, cloneEdges(prerequisites))
		want := canFinishByIndegree(numCourses, prerequisites)
		if got != want {
			t.Fatalf("canFinish(%d, %v) = %v, Kahn's algorithm = %v", numCourses, prerequisites, got, want)
		}
		if want {
			trueCount++
		} else {
			falseCount++
		}
	}

	if trueCount == 0 || falseCount == 0 {
		t.Fatalf("random inputs are lopsided (%d true, %d false), test is not covering both answers", trueCount, falseCount)
	}
}

func cloneEdges(edges [][]int) [][]int {
	out := make([][]int, len(edges))
	for i, e := range edges {
		out[i] = slices.Clone(e)
	}
	return out
}

// The input is only read.
func TestCanFinishDoesNotModifyInput(t *testing.T) {
	prerequisites := [][]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}}
	original := cloneEdges(prerequisites)

	canFinish(4, prerequisites)

	if len(prerequisites) != len(original) {
		t.Fatalf("canFinish changed the number of prerequisites: got %v, want %v", prerequisites, original)
	}
	for i := range original {
		if !slices.Equal(prerequisites[i], original[i]) {
			t.Errorf("canFinish modified its input: got %v, want %v", prerequisites, original)
			break
		}
	}
}

// Constraints allow 2000 courses and 5000 prerequisites. The ladder is the input
// that decides whether the DFS remembers courses it has already cleared: every
// course needs the next two, so the number of distinct paths from course 0 grows
// like Fibonacci. Without remembering, the DFS walks every one of those paths and
// never finishes. Runs in a goroutine so a blow-up fails at the budget instead
// of hanging the suite. Skipped with -short.
func TestCanFinishAtMaxSize(t *testing.T) {
	if testing.Short() {
		t.Skip("large inputs skipped in -short mode")
	}

	const n = 2000
	const budget = time.Second

	// every course i needs i+1 and i+2: 3997 edges, no cycle
	var ladder [][]int
	for i := range n {
		if i+1 < n {
			ladder = append(ladder, []int{i, i + 1})
		}
		if i+2 < n {
			ladder = append(ladder, []int{i, i + 2})
		}
	}

	// the same ladder with one edge back from the last course to the first
	ladderWithCycle := append(cloneEdges(ladder), []int{n - 1, 0})

	// a chain the full length of the course list: recursion 2000 deep
	var chain [][]int
	for i := 0; i+1 < n; i++ {
		chain = append(chain, []int{i, i + 1})
	}

	// random DAG at the edge limit: only ever point from a lower course to a
	// higher one, which cannot make a cycle
	rng := rand.New(rand.NewPCG(23, 24))
	seen := map[[2]int]bool{}
	var randomDAG [][]int
	for len(randomDAG) < 5000 {
		a, b := rng.IntN(n), rng.IntN(n)
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		if seen[[2]int{a, b}] {
			continue
		}
		seen[[2]int{a, b}] = true
		randomDAG = append(randomDAG, []int{a, b})
	}

	cases := []struct {
		name          string
		prerequisites [][]int
		want          bool
	}{
		{"ladder", ladder, true},
		{"ladder with a cycle", ladderWithCycle, false},
		{"full-length chain", chain, true},
		{"random DAG at the edge limit", randomDAG, true},
	}

	for _, c := range cases {
		done := make(chan bool, 1)
		start := time.Now()
		go func() { done <- canFinish(n, c.prerequisites) }()

		select {
		case got := <-done:
			if got != c.want {
				t.Errorf("%s (%d courses, %d prerequisites): got %v, want %v", c.name, n, len(c.prerequisites), got, c.want)
			}
			t.Logf("%s (%d courses, %d prerequisites): %v", c.name, n, len(c.prerequisites), time.Since(start).Round(time.Millisecond))
		case <-time.After(budget):
			t.Errorf("%s (%d courses, %d prerequisites): still running after %v, over budget", c.name, n, len(c.prerequisites), budget)
		}
	}
}
