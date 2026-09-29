package graph

import "testing"

// rottingByPerCellBFS is an independent oracle: for every fresh cell, walk
// outwards until a rotten one is found, and take the largest of those
// distances. Any fresh cell that never reaches a rotten one means -1. One
// search per cell, so it shares nothing with the single sweep under test.
func rottingByPerCellBFS(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])

	worst := 0
	for r := range rows {
		for c := range cols {
			if grid[r][c] != 1 {
				continue
			}

			seen := map[[2]int]bool{{r, c}: true}
			frontier := [][2]int{{r, c}}
			found := -1

			for steps := 0; len(frontier) > 0 && found < 0; steps++ {
				var next [][2]int
				for _, cell := range frontier {
					if grid[cell[0]][cell[1]] == 2 {
						found = steps
						break
					}
					for _, step := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
						nr, nc := cell[0]+step[0], cell[1]+step[1]
						if nr < 0 || nc < 0 || nr >= rows || nc >= cols {
							continue
						}
						if grid[nr][nc] == 0 || seen[[2]int{nr, nc}] {
							continue
						}
						seen[[2]int{nr, nc}] = true
						next = append(next, [2]int{nr, nc})
					}
				}
				frontier = next
			}

			if found < 0 {
				return -1
			}
			if found > worst {
				worst = found
			}
		}
	}

	return worst
}

func TestOrangesRottingExamples(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		// samples from the statement
		{"leetcode sample 1", [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}}, 4},
		{"leetcode sample 2", [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}}, -1},
		{"leetcode sample 3", [][]int{{0, 2}}, 0},
		{"neetcode sample 1", [][]int{{1, 1, 0}, {0, 1, 1}, {0, 1, 2}}, 4},
		{"neetcode sample 2", [][]int{{1, 0, 1}, {0, 2, 0}, {1, 0, 1}}, -1},

		// single cells
		{"one empty cell", [][]int{{0}}, 0},
		{"one rotten cell", [][]int{{2}}, 0},
		{"one fresh cell", [][]int{{1}}, -1},

		// nothing to do
		{"no fresh fruit", [][]int{{0, 2}, {2, 0}}, 0},
		{"all rotten already", [][]int{{2, 2}, {2, 2}}, 0},
		{"all empty", [][]int{{0, 0}, {0, 0}}, 0},
	}

	for _, c := range cases {
		if got := orangesRotting(cloneIntGrid(c.grid)); got != c.want {
			t.Errorf("%s: orangesRotting =\n%swant %d, got %d", c.name, renderDistances(c.grid), c.want, got)
		}
	}
}

// Fresh fruit that nothing can reach means -1, however much else rots fine.
func TestOrangesRottingUnreachableFresh(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
	}{
		{"fresh with no rotten anywhere", [][]int{{1, 1}, {1, 1}}},
		{"fresh cut off by empty cells", [][]int{{2, 0, 1}}},
		{"fresh behind a wall of empties", [][]int{{2, 1, 0, 1}, {2, 1, 0, 1}}},
		{"one corner unreachable", [][]int{{2, 1, 1}, {1, 1, 0}, {0, 0, 1}}},
	}

	for _, c := range cases {
		if got := orangesRotting(cloneIntGrid(c.grid)); got != -1 {
			t.Errorf("%s: orangesRotting = %d, want -1, for grid:\n%s", c.name, got, renderDistances(c.grid))
		}
	}
}

// The answer is the time the LAST fruit rots, and every rotten cell spreads at
// the same time. Starting from one rotten cell at a time gives bigger numbers
// on several of these.
func TestOrangesRottingSpreadsFromEverywhereAtOnce(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		{"rotten at both ends of a row", [][]int{{2, 1, 1, 1, 2}}, 2},
		{"rotten at one end of a row", [][]int{{2, 1, 1, 1, 1}}, 4},
		{"rotten in the middle", [][]int{{1, 1, 2, 1, 1}}, 2},
		{"rotten in opposite corners", [][]int{{2, 1, 1}, {1, 1, 1}, {1, 1, 2}}, 2},
		{"rotten column", [][]int{{2, 1, 1}, {2, 1, 1}, {2, 1, 1}}, 2},
		{"one fresh next to rotten", [][]int{{2, 1}}, 1},
		{"long snake around empties", [][]int{{2, 1, 1, 1}, {0, 0, 0, 1}, {1, 1, 1, 1}}, 8},
	}

	for _, c := range cases {
		if got := orangesRotting(cloneIntGrid(c.grid)); got != c.want {
			t.Errorf("%s: orangesRotting = %d, want %d, for grid:\n%s", c.name, got, c.want, renderDistances(c.grid))
		}
	}
}

// Cross-check against the one-search-per-cell oracle.
func TestOrangesRottingAgainstPerCellBFS(t *testing.T) {
	grids := [][][]int{
		{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}},
		{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}},
		{{1, 1, 0}, {0, 1, 1}, {0, 1, 2}},
		{{0, 2}},
		{{2, 2}, {1, 1}},
		{{1, 2, 1}, {2, 1, 2}, {1, 2, 1}},
		{{2, 0, 1, 1}, {1, 1, 1, 0}, {0, 1, 2, 1}},
		{{1, 1, 1}, {1, 1, 1}, {1, 1, 2}},
		{{0, 0, 0}, {0, 2, 0}, {0, 0, 0}},
		{{2, 1, 0, 2, 1}, {1, 0, 1, 2, 1}, {1, 0, 0, 2, 1}},
	}

	for i, grid := range grids {
		got, want := orangesRotting(cloneIntGrid(grid)), rottingByPerCellBFS(grid)
		if got != want {
			t.Errorf("grid %d: orangesRotting = %d, oracle = %d, for:\n%s", i, got, want, renderDistances(grid))
		}
	}
}

// With no empty cells and a single rotten corner, the last fruit to rot is the
// opposite corner, at its Manhattan distance away.
func TestOrangesRottingLarge(t *testing.T) {
	const rows, cols = 60, 40

	grid := make([][]int, rows)
	for r := range grid {
		grid[r] = make([]int, cols)
		for c := range grid[r] {
			grid[r][c] = 1
		}
	}
	grid[0][0] = 2

	if got, want := orangesRotting(grid), (rows-1)+(cols-1); got != want {
		t.Errorf("%dx%d grid rotting from one corner: got %d, want %d", rows, cols, got, want)
	}
}
