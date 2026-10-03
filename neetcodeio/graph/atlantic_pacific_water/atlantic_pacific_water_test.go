package atlanticpacificwater

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// reachesBothByFlowingDown is an independent oracle: from each cell, walk downhill
// (to neighbours with height <= current) and see which oceans we fall into.
// O((m*n)^2), so only for small grids.
func reachesBothByFlowingDown(heights [][]int) [][]int {
	rows, cols := len(heights), len(heights[0])
	out := [][]int{}
	for sr := range rows {
		for sc := range cols {
			seen := make([][]bool, rows)
			for i := range seen {
				seen[i] = make([]bool, cols)
			}
			pac, atl := false, false
			stack := [][2]int{{sr, sc}}
			seen[sr][sc] = true
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				r, c := cur[0], cur[1]
				if r == 0 || c == 0 {
					pac = true
				}
				if r == rows-1 || c == cols-1 {
					atl = true
				}
				for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
					nr, nc := r+d[0], c+d[1]
					if nr < 0 || nc < 0 || nr >= rows || nc >= cols || seen[nr][nc] {
						continue
					}
					if heights[nr][nc] <= heights[r][c] {
						seen[nr][nc] = true
						stack = append(stack, [2]int{nr, nc})
					}
				}
			}
			if pac && atl {
				out = append(out, []int{sr, sc})
			}
		}
	}
	return out
}

// sorted makes the comparison order-independent: the problem allows any order.
func sorted(cells [][]int) [][]int {
	out := slices.Clone(cells)
	slices.SortFunc(out, func(a, b []int) int {
		if a[0] != b[0] {
			return a[0] - b[0]
		}
		return a[1] - b[1]
	})
	return out
}

func equalCells(a, b [][]int) bool {
	return slices.EqualFunc(sorted(a), sorted(b), slices.Equal[[]int])
}

func TestPacificAtlantic(t *testing.T) {
	tests := []struct {
		name    string
		heights [][]int
		want    [][]int
	}{
		{
			name: "sample 1",
			heights: [][]int{
				{4, 2, 7, 3, 4},
				{7, 4, 6, 4, 7},
				{6, 3, 5, 3, 6},
			},
			want: [][]int{{0, 2}, {0, 4}, {1, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {2, 0}},
		},
		{
			name:    "sample 2: two cells, flat",
			heights: [][]int{{1}, {1}},
			want:    [][]int{{0, 0}, {1, 0}},
		},
		{
			name: "leetcode sample",
			heights: [][]int{
				{1, 2, 2, 3, 5},
				{3, 2, 3, 4, 4},
				{2, 4, 5, 3, 1},
				{6, 7, 1, 4, 5},
				{5, 1, 1, 2, 4},
			},
			want: [][]int{{0, 4}, {1, 3}, {1, 4}, {2, 2}, {3, 0}, {3, 1}, {4, 0}},
		},
		{
			name:    "single cell touches both oceans",
			heights: [][]int{{42}},
			want:    [][]int{{0, 0}},
		},
		{
			name:    "single row: every cell borders both",
			heights: [][]int{{5, 1, 9, 2}},
			want:    [][]int{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		},
		{
			name:    "single column: every cell borders both",
			heights: [][]int{{3}, {8}, {1}},
			want:    [][]int{{0, 0}, {1, 0}, {2, 0}},
		},
		{
			name:    "all flat: water moves everywhere",
			heights: [][]int{{2, 2, 2}, {2, 2, 2}, {2, 2, 2}},
			want:    [][]int{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			name: "pit in the middle reaches nothing",
			heights: [][]int{
				{5, 5, 5},
				{5, 0, 5},
				{5, 5, 5},
			},
			want: [][]int{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			name: "peak in the middle reaches both",
			heights: [][]int{
				{1, 1, 1},
				{1, 9, 1},
				{1, 1, 1},
			},
			want: [][]int{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			name: "near miss: strictly rising toward atlantic corner",
			heights: [][]int{
				{1, 2},
				{2, 3},
			},
			// (0,0) is Pacific-only: it can't climb to the Atlantic.
			want: [][]int{{0, 1}, {1, 0}, {1, 1}},
		},
		{
			name: "near miss: wall blocks pacific from bottom-right basin",
			heights: [][]int{
				{1, 1, 1, 1},
				{1, 9, 9, 9},
				{1, 9, 0, 0},
				{1, 9, 0, 0},
			},
			// The basin (2..3, 2..3) is Atlantic-only: the 9s wall it off from the Pacific.
			want: [][]int{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 0}, {1, 1}, {1, 2}, {1, 3}, {2, 0}, {2, 1}, {3, 0}, {3, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pacificAtlantic(tt.heights)
			if !equalCells(got, tt.want) {
				t.Errorf("got %v, want %v", sorted(got), sorted(tt.want))
			}
			if oracle := reachesBothByFlowingDown(tt.heights); !equalCells(oracle, tt.want) {
				t.Fatalf("test table is wrong: oracle says %v", sorted(oracle))
			}
		})
	}
}

func TestPacificAtlanticRandomAgainstOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := range 500 {
		rows, cols := 1+rng.IntN(8), 1+rng.IntN(8)
		maxH := 1 + rng.IntN(6) // small range so ties (equal heights) are common
		heights := make([][]int, rows)
		for r := range heights {
			heights[r] = make([]int, cols)
			for c := range heights[r] {
				heights[r][c] = rng.IntN(maxH)
			}
		}
		got := pacificAtlantic(heights)
		want := reachesBothByFlowingDown(heights)
		if !equalCells(got, want) {
			t.Fatalf("case %d heights=%v\ngot  %v\nwant %v", i, heights, sorted(got), sorted(want))
		}
	}
}

// Max size: 200x200 with a snake-shaped ramp so the DFS goes 40,000 deep.
func TestPacificAtlanticMaxSize(t *testing.T) {
	if testing.Short() {
		t.Skip("max-size test")
	}
	const n = 200
	heights := make([][]int, n)
	for r := range heights {
		heights[r] = make([]int, n)
		for c := range heights[r] {
			if r%2 == 0 {
				heights[r][c] = r*n + c
			} else {
				heights[r][c] = r*n + (n - 1 - c)
			}
		}
	}
	got := pacificAtlantic(heights)
	if len(got) == 0 {
		t.Fatalf("expected at least the top-right / bottom-left style cells, got none")
	}
}
