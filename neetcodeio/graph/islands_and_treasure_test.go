package graph

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// INF marks a land cell that has not been reached yet, as in the statement.
const INF = math.MaxInt32

// distancesByBFSPerCell is an independent oracle: for every land cell, walk
// outwards level by level until a treasure is found, and record how many steps
// it took. One search per cell, so it shares nothing with the single sweep under
// test. O((m*n)^2), fine for small grids.
func distancesByBFSPerCell(grid [][]int) [][]int {
	rows, cols := len(grid), len(grid[0])

	out := make([][]int, rows)
	for r := range out {
		out[r] = make([]int, cols)
		copy(out[r], grid[r])
	}

	for r := range rows {
		for c := range cols {
			if grid[r][c] != INF {
				continue
			}

			seen := map[[2]int]bool{{r, c}: true}
			frontier := [][2]int{{r, c}}

			for steps := 0; len(frontier) > 0; steps++ {
				var next [][2]int
				for _, cell := range frontier {
					if grid[cell[0]][cell[1]] == 0 {
						out[r][c] = steps
						next = nil
						frontier = nil
						break
					}
					for _, step := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
						nr, nc := cell[0]+step[0], cell[1]+step[1]
						if nr < 0 || nc < 0 || nr >= rows || nc >= cols {
							continue
						}
						if grid[nr][nc] == -1 || seen[[2]int{nr, nc}] {
							continue
						}
						seen[[2]int{nr, nc}] = true
						next = append(next, [2]int{nr, nc})
					}
				}
				if frontier == nil {
					break
				}
				frontier = next
			}
		}
	}

	return out
}

func cloneIntGrid(grid [][]int) [][]int {
	out := make([][]int, len(grid))
	for i, row := range grid {
		out[i] = make([]int, len(row))
		copy(out[i], row)
	}
	return out
}

func renderDistances(grid [][]int) string {
	var b strings.Builder
	for _, row := range grid {
		for i, v := range row {
			if i > 0 {
				b.WriteByte(' ')
			}
			switch v {
			case INF:
				b.WriteString("INF")
			case -1:
				b.WriteString(" -1")
			default:
				fmt.Fprintf(&b, "%3d", v)
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func gridsEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestIslandsAndTreasureExamples(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want [][]int
	}{
		{
			name: "statement sample",
			grid: [][]int{
				{INF, -1, 0, INF},
				{INF, INF, INF, -1},
				{INF, -1, INF, -1},
				{0, -1, INF, INF},
			},
			want: [][]int{
				{3, -1, 0, 1},
				{2, 2, 1, -1},
				{1, -1, 2, -1},
				{0, -1, 3, 4},
			},
		},
		{
			name: "one treasure in a corner",
			grid: [][]int{
				{0, -1},
				{INF, INF},
			},
			want: [][]int{
				{0, -1},
				{1, 2},
			},
		},
		{
			name: "single treasure",
			grid: [][]int{{0}},
			want: [][]int{{0}},
		},
		{
			name: "single water cell",
			grid: [][]int{{-1}},
			want: [][]int{{-1}},
		},
		{
			name: "single unreachable land cell",
			grid: [][]int{{INF}},
			want: [][]int{{INF}},
		},
	}

	for _, c := range cases {
		got := cloneIntGrid(c.grid)
		islandsAndTreasure(got)

		if !gridsEqual(got, c.want) {
			t.Errorf("%s:\ngot:\n%swant:\n%s", c.name, renderDistances(got), renderDistances(c.want))
		}
	}
}

// Cells that no treasure can reach keep INF, and water is never overwritten.
func TestIslandsAndTreasureUnreachable(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want [][]int
	}{
		{
			name: "no treasure at all",
			grid: [][]int{
				{INF, INF},
				{INF, -1},
			},
			want: [][]int{
				{INF, INF},
				{INF, -1},
			},
		},
		{
			name: "land sealed off by water",
			grid: [][]int{
				{0, -1, INF},
				{INF, -1, INF},
				{INF, -1, INF},
			},
			want: [][]int{
				{0, -1, INF},
				{1, -1, INF},
				{2, -1, INF},
			},
		},
		{
			name: "only water and treasure",
			grid: [][]int{
				{0, -1},
				{-1, 0},
			},
			want: [][]int{
				{0, -1},
				{-1, 0},
			},
		},
		{
			name: "every cell is a treasure",
			grid: [][]int{
				{0, 0},
				{0, 0},
			},
			want: [][]int{
				{0, 0},
				{0, 0},
			},
		},
	}

	for _, c := range cases {
		got := cloneIntGrid(c.grid)
		islandsAndTreasure(got)

		if !gridsEqual(got, c.want) {
			t.Errorf("%s:\ngot:\n%swant:\n%s", c.name, renderDistances(got), renderDistances(c.want))
		}
	}
}

// With several treasures, every cell must take the nearest one. A sweep that
// starts from one treasure at a time, or that stops at the first treasure it
// finds, gets these wrong.
func TestIslandsAndTreasureNearestWins(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want [][]int
	}{
		{
			name: "treasure at both ends of a row",
			grid: [][]int{{0, INF, INF, INF, 0}},
			want: [][]int{{0, 1, 2, 1, 0}},
		},
		{
			name: "treasure at both ends of a column",
			grid: [][]int{{0}, {INF}, {INF}, {0}},
			want: [][]int{{0}, {1}, {1}, {0}},
		},
		{
			name: "treasures in opposite corners",
			grid: [][]int{
				{0, INF, INF},
				{INF, INF, INF},
				{INF, INF, 0},
			},
			want: [][]int{
				{0, 1, 2},
				{1, 2, 1},
				{2, 1, 0},
			},
		},
		{
			name: "a wall makes the far treasure the nearest one",
			grid: [][]int{
				{0, -1, INF},
				{-1, -1, INF},
				{INF, INF, 0},
			},
			want: [][]int{
				{0, -1, 2},
				{-1, -1, 1},
				{2, 1, 0},
			},
		},
	}

	for _, c := range cases {
		got := cloneIntGrid(c.grid)
		islandsAndTreasure(got)

		if !gridsEqual(got, c.want) {
			t.Errorf("%s:\ngot:\n%swant:\n%s", c.name, renderDistances(got), renderDistances(c.want))
		}
	}
}

// Cross-check against the one-search-per-cell oracle.
func TestIslandsAndTreasureAgainstPerCellBFS(t *testing.T) {
	grids := [][][]int{
		{{INF, -1, 0, INF}, {INF, INF, INF, -1}, {INF, -1, INF, -1}, {0, -1, INF, INF}},
		{{0, INF, INF}, {INF, INF, INF}, {INF, INF, INF}},
		{{INF, INF, INF}, {INF, 0, INF}, {INF, INF, INF}},
		{{0, -1, INF, INF}, {INF, INF, INF, -1}, {-1, -1, 0, INF}},
		{{INF, INF}, {INF, INF}},
		{{0, 0}, {INF, INF}},
		{{-1, INF, -1}, {INF, 0, INF}, {-1, INF, -1}},
		{{INF, INF, INF, INF, INF}, {INF, -1, -1, -1, INF}, {INF, -1, 0, -1, INF}, {INF, -1, -1, -1, INF}, {INF, INF, INF, INF, INF}},
	}

	for i, grid := range grids {
		got := cloneIntGrid(grid)
		islandsAndTreasure(got)

		want := distancesByBFSPerCell(grid)

		if !gridsEqual(got, want) {
			t.Errorf("grid %d:\ngot:\n%soracle:\n%s", i, renderDistances(got), renderDistances(want))
		}
	}
}

// With no walls and a single treasure in a corner, the answer for every cell is
// its Manhattan distance to that corner. Constraints allow a few hundred cells
// per side.
func TestIslandsAndTreasureLarge(t *testing.T) {
	const n = 100

	grid := make([][]int, n)
	for r := range grid {
		grid[r] = make([]int, n)
		for c := range grid[r] {
			grid[r][c] = INF
		}
	}
	grid[0][0] = 0

	islandsAndTreasure(grid)

	for r := range n {
		for c := range n {
			if want := r + c; grid[r][c] != want {
				t.Fatalf("cell (%d,%d) = %d, want %d", r, c, grid[r][c], want)
			}
		}
	}
}
