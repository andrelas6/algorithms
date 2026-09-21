package graph

import (
	"strings"
	"testing"
)

// intGridOf builds the grid from one string per row, '0' for open and '1' for
// blocked, so the tests read like the pictures in the statement.
func intGridOf(rows ...string) [][]int {
	grid := make([][]int, len(rows))
	for i, row := range rows {
		grid[i] = make([]int, len(row))
		for j := range row {
			if row[j] == '1' {
				grid[i][j] = 1
			}
		}
	}
	return grid
}

func renderInts(grid [][]int) string {
	var b strings.Builder
	for _, row := range grid {
		for _, v := range row {
			if v == 1 {
				b.WriteByte('1')
			} else {
				b.WriteByte('0')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// countPathsByBitmask is an independent oracle: an iterative walk with an
// explicit stack, carrying the set of used cells as a bitmask instead of a map,
// and no recursion. Limited to grids of at most 63 cells.
func countPathsByBitmask(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	if grid[0][0] == 1 {
		return 0
	}

	type state struct {
		r, c int
		used uint64
	}

	bit := func(r, c int) uint64 { return 1 << (r*cols + c) }

	stack := []state{{0, 0, bit(0, 0)}}
	count := 0

	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if s.r == rows-1 && s.c == cols-1 {
			count++
			continue
		}

		steps := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, step := range steps {
			nr, nc := s.r+step[0], s.c+step[1]
			if nr < 0 || nc < 0 || nr >= rows || nc >= cols {
				continue
			}
			if grid[nr][nc] == 1 || s.used&bit(nr, nc) != 0 {
				continue
			}
			stack = append(stack, state{nr, nc, s.used | bit(nr, nc)})
		}
	}

	return count
}

// Counts worked out by hand, and for the open squares the known counts of
// corner-to-corner self-avoiding walks: 1, 2, 12, 184 for sizes 1 to 4.
func TestCountPathsByHand(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		{"single open cell", intGridOf("0"), 1},
		{"single blocked cell", intGridOf("1"), 0},

		{"one row, open", intGridOf("00"), 1},
		{"one row, blocked", intGridOf("01"), 0},
		{"one row of four", intGridOf("0000"), 1},
		{"one row with a wall", intGridOf("0010"), 0},
		{"one column", intGridOf("0", "0", "0"), 1},
		{"one column with a wall", intGridOf("0", "1", "0"), 0},

		{"two by two, open", intGridOf("00", "00"), 2},
		{"two by two, one wall", intGridOf("01", "00"), 1},
		{"two by two, both middles walled", intGridOf("01", "10"), 0},
		{"two by two, end blocked", intGridOf("00", "01"), 0},
		{"two by two, start blocked", intGridOf("10", "00"), 0},

		{"three by three, open", intGridOf("000", "000", "000"), 12},
		{"four by four, open", intGridOf("0000", "0000", "0000", "0000"), 184},

		{"only one way through", intGridOf("001", "101", "000"), 1},
		{"walled off", intGridOf("001", "110", "000"), 0},
	}

	for _, c := range cases {
		if got := countPaths(c.grid); got != c.want {
			t.Errorf("%s: countPaths = %d, want %d, for grid:\n%s", c.name, got, c.want, renderInts(c.grid))
		}
	}
}

// Cross-check against the bitmask oracle on a spread of grids.
func TestCountPathsAgainstBitmask(t *testing.T) {
	grids := [][][]int{
		intGridOf("0"),
		intGridOf("00"),
		intGridOf("00", "00"),
		intGridOf("000", "000", "000"),
		intGridOf("0000", "1100", "0001", "0100"),
		intGridOf("0000", "0000", "0000"),
		intGridOf("00000", "01110", "00000"),
		intGridOf("010", "000", "010"),
		intGridOf("0100", "0010", "0000", "0110"),
		intGridOf("00000", "00000", "00000", "00000"),
		intGridOf("0010", "0100", "0001", "0000"),
	}

	for i, grid := range grids {
		got, want := countPaths(grid), countPathsByBitmask(grid)
		if got != want {
			t.Errorf("grid %d: countPaths = %d, bitmask oracle = %d, for:\n%s", i, got, want, renderInts(grid))
		}
	}
}

// Walls placed one at a time: each wall can only ever remove paths, never add
// them, and walling the start or the end always gives 0.
func TestCountPathsWallsOnlyRemovePaths(t *testing.T) {
	base := intGridOf("0000", "0000", "0000", "0000")
	open := countPaths(base)

	rows, cols := len(base), len(base[0])
	for r := range rows {
		for c := range cols {
			walled := intGridOf("0000", "0000", "0000", "0000")
			walled[r][c] = 1

			got := countPaths(walled)
			if got > open {
				t.Errorf("wall at (%d,%d): %d paths, more than the %d with no walls", r, c, got, open)
			}
			if (r == 0 && c == 0) || (r == rows-1 && c == cols-1) {
				if got != 0 {
					t.Errorf("wall at (%d,%d) blocks the start or the end: got %d, want 0", r, c, got)
				}
			}
		}
	}
}

// Counting is read-only.
func TestCountPathsDoesNotModifyGrid(t *testing.T) {
	grid := intGridOf("0000", "1100", "0001", "0100")
	before := renderInts(grid)

	countPaths(grid)

	if after := renderInts(grid); after != before {
		t.Errorf("countPaths modified the grid:\nbefore:\n%safter:\n%s", before, after)
	}
}

// The search undoes its own bookkeeping: every cell it marks on the way down is
// unmarked on the way back up, so the set comes back empty.
func TestDfsCountPathsLeavesNothingBehind(t *testing.T) {
	grids := [][][]int{
		intGridOf("000", "000", "000"),
		intGridOf("0000", "1100", "0001", "0100"),
		intGridOf("01", "00"),
	}

	for i, grid := range grids {
		v := make(visited)

		dfsCountPaths(grid, 0, 0, v)

		if len(v) != 0 {
			t.Errorf("grid %d: %d positions left marked after the search finished", i, len(v))
		}
	}
}

// A 5x5 open grid has 8512 corner-to-corner paths. This is the size where the
// exponential shape of the search starts to show.
func TestCountPathsFiveByFive(t *testing.T) {
	if testing.Short() {
		t.Skip("large inputs skipped in -short mode")
	}

	grid := intGridOf("00000", "00000", "00000", "00000", "00000")

	if got := countPaths(grid); got != 8512 {
		t.Errorf("five by five open grid: countPaths = %d, want 8512", got)
	}
}
