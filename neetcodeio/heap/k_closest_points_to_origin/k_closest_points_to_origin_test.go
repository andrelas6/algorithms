package kclosestpointstoorigin

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"testing"
)

// sortPoints puts points in a fixed order so two answers can be compared
// regardless of the order the heap handed them back in.
func sortPoints(points [][]int) [][]int {
	out := make([][]int, len(points))
	for i, p := range points {
		out[i] = slices.Clone(p)
	}
	slices.SortFunc(out, func(a, b []int) int {
		if c := cmp.Compare(a[0], b[0]); c != 0 {
			return c
		}
		return cmp.Compare(a[1], b[1])
	})
	return out
}

func equalPoints(a, b [][]int) bool {
	return slices.EqualFunc(sortPoints(a), sortPoints(b), slices.Equal[[]int])
}

// oracle sorts every point by squared distance and takes the first k.
// No heap involved, so it can't share a heap bug.
func oracle(points [][]int, k int) [][]int {
	sorted := sortPoints(points)
	slices.SortStableFunc(sorted, func(a, b []int) int {
		return cmp.Compare(a[0]*a[0]+a[1]*a[1], b[0]*b[0]+b[1]*b[1])
	})
	return sorted[:k]
}

func TestKClosest(t *testing.T) {
	cases := []struct {
		name   string
		points [][]int
		k      int
		want   [][]int
	}{
		{"neetcode sample", [][]int{{0, 2}, {2, 2}}, 1, [][]int{{0, 2}}},
		{"leetcode sample 1", [][]int{{1, 3}, {-2, 2}}, 1, [][]int{{-2, 2}}},
		{"leetcode sample 2", [][]int{{3, 3}, {5, -1}, {-2, 4}}, 2, [][]int{{3, 3}, {-2, 4}}},
		{"single point", [][]int{{7, -7}}, 1, [][]int{{7, -7}}},
		{"k equals n returns everything", [][]int{{1, 1}, {2, 2}, {3, 3}}, 3, [][]int{{1, 1}, {2, 2}, {3, 3}}},
		{"origin itself is closest", [][]int{{5, 5}, {0, 0}, {1, 0}}, 1, [][]int{{0, 0}}},
		{"negatives square to positive", [][]int{{-1, -1}, {2, 0}, {0, -3}}, 2, [][]int{{-1, -1}, {2, 0}}},
		// sum of coords would say (3,-3) is closest at 0; squared distance says (1,1)
		{"sum of coords is not distance", [][]int{{3, -3}, {1, 1}}, 1, [][]int{{1, 1}}},
		// |x|+|y| would call these equal (4 vs 4); squared says 8 vs 16
		{"manhattan is not euclidean", [][]int{{4, 0}, {2, 2}}, 1, [][]int{{2, 2}}},
		{"duplicate points", [][]int{{1, 1}, {1, 1}, {5, 5}}, 2, [][]int{{1, 1}, {1, 1}}},
		{"max constraint coords don't overflow", [][]int{{10000, 10000}, {-10000, 9999}}, 1, [][]int{{-10000, 9999}}},
		{"input already sorted far to near", [][]int{{9, 9}, {5, 5}, {3, 3}, {1, 1}}, 2, [][]int{{3, 3}, {1, 1}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := sortPoints(tc.points) // copy, so we can check the input isn't mangled
			got := kClosest(input, tc.k)
			if len(got) != tc.k {
				t.Fatalf("kClosest(%v, %d) returned %d points, want %d", tc.points, tc.k, len(got), tc.k)
			}
			if !equalPoints(got, tc.want) {
				t.Errorf("kClosest(%v, %d) = %v, want %v (any order)", tc.points, tc.k, got, tc.want)
			}
		})
	}
}

func TestKClosestAgainstOracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for iter := range 500 {
		n := 1 + r.IntN(50)
		// small coordinate range forces lots of distance ties
		span := 1 + r.IntN(10)
		points := make([][]int, n)
		for i := range points {
			points[i] = []int{r.IntN(2*span+1) - span, r.IntN(2*span+1) - span}
		}
		k := 1 + r.IntN(n)

		got := kClosest(points, k)
		want := oracle(points, k)

		// Ties at the k-th distance mean more than one valid answer, so compare
		// the multiset of distances, not the exact points.
		dists := func(ps [][]int) []int {
			d := make([]int, len(ps))
			for i, p := range ps {
				d[i] = p[0]*p[0] + p[1]*p[1]
			}
			slices.Sort(d)
			return d
		}
		if !slices.Equal(dists(got), dists(want)) {
			t.Fatalf("iter %d: kClosest(%v, %d) distances = %v, want %v", iter, points, k, dists(got), dists(want))
		}
	}
}
