package heap

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

// streamOracle is an independent reference: keep every value seen in a sorted
// slice and read off the kth largest directly. No heap, so a bug in the heap
// bookkeeping cannot hide behind the same mistake.
type streamOracle struct {
	k      int
	values []int
}

func newStreamOracle(k int, nums []int) *streamOracle {
	o := &streamOracle{k: k, values: slices.Clone(nums)}
	slices.Sort(o.values)
	return o
}

// add records the value and reports the kth largest so far. ok is false while
// fewer than k values have arrived, when there is no kth largest yet.
func (o *streamOracle) add(val int) (kth int, ok bool) {
	pos, _ := slices.BinarySearch(o.values, val)
	o.values = slices.Insert(o.values, pos, val)

	if len(o.values) < o.k {
		return 0, false
	}
	return o.values[len(o.values)-o.k], true
}

func TestKthLargestStreamExamples(t *testing.T) {
	t.Run("leetcode sample", func(t *testing.T) {
		kth := Constructor(3, []int{4, 5, 8, 2})

		cases := []struct{ add, want int }{
			{3, 4},
			{5, 5},
			{10, 5},
			{9, 8},
			{4, 8},
		}
		for _, c := range cases {
			if got := kth.Add(c.add); got != c.want {
				t.Errorf("Add(%d) = %d, want %d", c.add, got, c.want)
			}
		}
	})

	t.Run("neetcode sample", func(t *testing.T) {
		kth := Constructor(3, []int{1, 2, 3, 3})

		cases := []struct{ add, want int }{
			{3, 3},
			{5, 3},
			{6, 3},
			{7, 5},
			{8, 6},
		}
		for _, c := range cases {
			if got := kth.Add(c.add); got != c.want {
				t.Errorf("Add(%d) = %d, want %d", c.add, got, c.want)
			}
		}
	})
}

func TestKthLargestStreamSmallCases(t *testing.T) {
	t.Run("k of 1 tracks the maximum", func(t *testing.T) {
		kth := Constructor(1, []int{})

		cases := []struct{ add, want int }{
			{5, 5},
			{3, 5},
			{9, 9},
			{9, 9},
			{-2, 9},
		}
		for _, c := range cases {
			if got := kth.Add(c.add); got != c.want {
				t.Errorf("Add(%d) = %d, want %d", c.add, got, c.want)
			}
		}
	})

	t.Run("k equal to the starting length", func(t *testing.T) {
		kth := Constructor(2, []int{7, 9})

		if got := kth.Add(8); got != 8 {
			t.Errorf("Add(8) = %d, want 8", got)
		}
		if got := kth.Add(1); got != 8 {
			t.Errorf("Add(1) = %d, want 8", got)
		}
	})

	t.Run("starting shorter than k", func(t *testing.T) {
		kth := Constructor(3, []int{4, 5})

		if got := kth.Add(10); got != 4 {
			t.Errorf("Add(10) = %d, want 4", got)
		}
		if got := kth.Add(6); got != 5 {
			t.Errorf("Add(6) = %d, want 5", got)
		}
	})

	t.Run("duplicates count separately", func(t *testing.T) {
		kth := Constructor(2, []int{5, 5})

		if got := kth.Add(5); got != 5 {
			t.Errorf("Add(5) = %d, want 5", got)
		}
		if got := kth.Add(4); got != 5 {
			t.Errorf("Add(4) = %d, want 5", got)
		}
	})

	t.Run("negatives", func(t *testing.T) {
		kth := Constructor(2, []int{-10, -5})

		if got := kth.Add(-7); got != -7 {
			t.Errorf("Add(-7) = %d, want -7", got)
		}
		if got := kth.Add(-1); got != -5 {
			t.Errorf("Add(-1) = %d, want -5", got)
		}
	})
}

// Values arriving in the worst orders: only ever bigger, only ever smaller, and
// all the same.
func TestKthLargestStreamOrders(t *testing.T) {
	t.Run("increasing", func(t *testing.T) {
		kth := Constructor(3, []int{1, 2, 3})
		oracle := newStreamOracle(3, []int{1, 2, 3})

		for v := 4; v <= 20; v++ {
			got := kth.Add(v)
			if want, ok := oracle.add(v); ok && got != want {
				t.Fatalf("Add(%d) = %d, want %d", v, got, want)
			}
		}
	})

	t.Run("decreasing", func(t *testing.T) {
		kth := Constructor(3, []int{20, 19, 18})
		oracle := newStreamOracle(3, []int{20, 19, 18})

		for v := 17; v >= 1; v-- {
			got := kth.Add(v)
			if want, ok := oracle.add(v); ok && got != want {
				t.Fatalf("Add(%d) = %d, want %d", v, got, want)
			}
		}
	})

	t.Run("all the same value", func(t *testing.T) {
		kth := Constructor(4, []int{2, 2, 2, 2})

		for range 10 {
			if got := kth.Add(2); got != 2 {
				t.Fatalf("Add(2) = %d, want 2", got)
			}
		}
	})
}

// Random streams against the sorted-slice reference, for several values of k.
func TestKthLargestStreamAgainstOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(35, 36))

	for range 200 {
		k := 1 + rng.IntN(5)

		start := make([]int, rng.IntN(6))
		for i := range start {
			start[i] = rng.IntN(41) - 20
		}

		kth := Constructor(k, slices.Clone(start))
		oracle := newStreamOracle(k, start)

		for step := range 40 {
			v := rng.IntN(41) - 20

			got := kth.Add(v)

			// the kth largest only exists once k values have arrived
			if want, ok := oracle.add(v); ok && got != want {
				t.Fatalf("k=%d start=%v step %d: Add(%d) = %d, want %d", k, start, step, v, got, want)
			}
		}
	}
}

// Each instance keeps its own values.
func TestKthLargestStreamInstancesAreIndependent(t *testing.T) {
	a := Constructor(2, []int{1, 2})
	b := Constructor(2, []int{100, 200})

	if got := a.Add(3); got != 2 {
		t.Errorf("first instance: Add(3) = %d, want 2", got)
	}
	if got := b.Add(300); got != 200 {
		t.Errorf("second instance: Add(300) = %d, want 200", got)
	}
	if got := a.Add(4); got != 3 {
		t.Errorf("first instance after the second was used: Add(4) = %d, want 3", got)
	}
}

// The point of the heap is that it stays the size of k rather than the size of
// the stream. Constraints allow 10^4 calls.
func TestKthLargestStreamLong(t *testing.T) {
	if testing.Short() {
		t.Skip("large inputs skipped in -short mode")
	}

	const n = 10_000
	const k = 10
	const budget = time.Second

	rng := rand.New(rand.NewPCG(37, 38))

	kth := Constructor(k, []int{})
	oracle := newStreamOracle(k, []int{})

	start := time.Now()
	for step := range n {
		v := rng.IntN(20_001) - 10_000

		got := kth.Add(v)

		if want, ok := oracle.add(v); ok && got != want {
			t.Fatalf("step %d: Add(%d) = %d, want %d", step, v, got, want)
		}
	}
	elapsed := time.Since(start)

	if elapsed > budget {
		t.Errorf("%d adds took %v, over the %v budget", n, elapsed.Round(time.Millisecond), budget)
	}
	if size := kth.minHeap.Len(); size != k {
		t.Errorf("heap holds %d values after %d adds, want %d", size, n, k)
	}
}
