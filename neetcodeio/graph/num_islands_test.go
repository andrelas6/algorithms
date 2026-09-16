package graph

import (
	"math/rand/v2"
	"strings"
	"testing"
	"time"
)

// gridOf builds the grid from one string per row, so the tests read like the
// pictures in the statement.
func gridOf(rows ...string) [][]byte {
	grid := make([][]byte, len(rows))
	for i, row := range rows {
		grid[i] = []byte(row)
	}
	return grid
}

func cloneGrid(grid [][]byte) [][]byte {
	out := make([][]byte, len(grid))
	for i, row := range grid {
		out[i] = append([]byte(nil), row...)
	}
	return out
}

// countIslandsByFlooding is an independent oracle: a recursive flood fill that
// sinks each island it finds by overwriting the grid copy. No queue, no visited
// map, so a bug in the BFS bookkeeping cannot hide behind the same mistake.
func countIslandsByFlooding(grid [][]byte) int {
	work := cloneGrid(grid)

	var sink func(r, c int)
	sink = func(r, c int) {
		if r < 0 || c < 0 || r >= len(work) || c >= len(work[r]) || work[r][c] != '1' {
			return
		}
		work[r][c] = '0'
		sink(r-1, c)
		sink(r+1, c)
		sink(r, c-1)
		sink(r, c+1)
	}

	count := 0
	for r := range work {
		for c := range work[r] {
			if work[r][c] == '1' {
				count++
				sink(r, c)
			}
		}
	}
	return count
}

func TestNumIslandsExamples(t *testing.T) {
	cases := []struct {
		name string
		grid [][]byte
		want int
	}{
		{"neetcode sample 1", gridOf(
			"01110",
			"01010",
			"11000",
			"00000",
		), 1},
		{"neetcode sample 2", gridOf(
			"11001",
			"11001",
			"00100",
			"00011",
		), 4},
		{"leetcode sample 1", gridOf(
			"11110",
			"11010",
			"11000",
			"00000",
		), 1},
		{"leetcode sample 2", gridOf(
			"11000",
			"11000",
			"00100",
			"00011",
		), 3},
	}

	for _, c := range cases {
		if got := numIslands(c.grid); got != c.want {
			t.Errorf("%s: numIslands = %d, want %d", c.name, got, c.want)
		}
	}
}

func TestNumIslandsSmallGrids(t *testing.T) {
	cases := []struct {
		name string
		grid [][]byte
		want int
	}{
		// single cell
		{"one land cell", gridOf("1"), 1},
		{"one water cell", gridOf("0"), 0},

		// single row: neighbours only to the left and right
		{"row, one island", gridOf("111"), 1},
		{"row, two islands", gridOf("101"), 2},
		{"row, three islands", gridOf("10101"), 3},
		{"row, all water", gridOf("000"), 0},

		// single column: neighbours only up and down
		{"column, one island", gridOf("1", "1", "1"), 1},
		{"column, two islands", gridOf("1", "0", "1"), 2},

		// whole grid
		{"all land", gridOf("111", "111", "111"), 1},
		{"all water", gridOf("000", "000", "000"), 0},
	}

	for _, c := range cases {
		if got := numIslands(c.grid); got != c.want {
			t.Errorf("%s: numIslands = %d, want %d", c.name, got, c.want)
		}
	}
}

// The trap: only up, down, left and right connect. Cells touching at a corner
// are separate islands. A solution that also walks the diagonals returns 1 for
// every one of these.
func TestNumIslandsDiagonalsDoNotConnect(t *testing.T) {
	cases := []struct {
		name string
		grid [][]byte
		want int
	}{
		{"two cells on a diagonal", gridOf(
			"10",
			"01",
		), 2},
		{"anti-diagonal", gridOf(
			"01",
			"10",
		), 2},
		{"checkerboard 3x3", gridOf(
			"101",
			"010",
			"101",
		), 5},
		{"checkerboard 4x4", gridOf(
			"1010",
			"0101",
			"1010",
			"0101",
		), 8},
		{"diagonal chain", gridOf(
			"1000",
			"0100",
			"0010",
			"0001",
		), 4},
	}

	for _, c := range cases {
		if got := numIslands(c.grid); got != c.want {
			t.Errorf("%s: numIslands = %d, want %d", c.name, got, c.want)
		}
	}
}

// Shapes that force the search to turn corners and come back on itself: a
// solution that only walks in one direction, or that stops at the first branch,
// splits these into several islands.
func TestNumIslandsAwkwardShapes(t *testing.T) {
	cases := []struct {
		name string
		grid [][]byte
		want int
	}{
		{"snake", gridOf(
			"11111",
			"00001",
			"11111",
			"10000",
			"11111",
		), 1},
		{"ring with a hole", gridOf(
			"11111",
			"10001",
			"10001",
			"10001",
			"11111",
		), 1},
		{"island inside a ring", gridOf(
			"11111",
			"10001",
			"10101",
			"10001",
			"11111",
		), 2},
		{"plus sign", gridOf(
			"010",
			"111",
			"010",
		), 1},
		{"u shape", gridOf(
			"101",
			"101",
			"111",
		), 1},
		{"four corners", gridOf(
			"101",
			"000",
			"101",
		), 4},
		{"touching only at the edges", gridOf(
			"1100",
			"1100",
			"0011",
			"0011",
		), 2},
		{"comb", gridOf(
			"10101",
			"10101",
			"11111",
		), 1},
	}

	for _, c := range cases {
		if got := numIslands(c.grid); got != c.want {
			t.Errorf("%s: numIslands = %d, want %d", c.name, got, c.want)
		}
	}
}

// Random grids against the flood-fill oracle. Different land densities: sparse
// grids make many one-cell islands, dense ones make a few big tangled shapes.
func TestNumIslandsAgainstFloodFill(t *testing.T) {
	rng := rand.New(rand.NewPCG(15, 16))

	for range 1500 {
		rows := 1 + rng.IntN(8)
		cols := 1 + rng.IntN(8)
		density := []int{20, 45, 70, 90}[rng.IntN(4)]

		grid := make([][]byte, rows)
		for r := range grid {
			grid[r] = make([]byte, cols)
			for c := range grid[r] {
				if rng.IntN(100) < density {
					grid[r][c] = '1'
				} else {
					grid[r][c] = '0'
				}
			}
		}

		got, want := numIslands(cloneGrid(grid)), countIslandsByFlooding(grid)
		if got != want {
			t.Fatalf("numIslands = %d, flood fill = %d, for grid:\n%s", got, want, render(grid))
		}
	}
}

func render(grid [][]byte) string {
	var out strings.Builder
	for _, row := range grid {
		out.Write(row)
		out.WriteByte('\n')
	}
	return out.String()
}

// Counting islands should not rewrite the caller's grid. The usual alternative
// solution sinks each island in place, so this pins which of the two this is.
func TestNumIslandsDoesNotModifyGrid(t *testing.T) {
	grid := gridOf(
		"11000",
		"11000",
		"00100",
		"00011",
	)
	before := render(grid)

	numIslands(grid)

	if after := render(grid); after != before {
		t.Errorf("numIslands modified the grid:\nbefore:\n%safter:\n%s", before, after)
	}
}

// Constraints allow 300x300. Both extremes: one solid island, and the
// checkerboard, which is the maximum number of islands a grid can hold.
func TestNumIslandsAtMaxSize(t *testing.T) {
	const n = 300
	const budget = 2 * time.Second

	allLand := make([][]byte, n)
	checker := make([][]byte, n)
	for r := range n {
		allLand[r] = make([]byte, n)
		checker[r] = make([]byte, n)
		for c := range n {
			allLand[r][c] = '1'
			if (r+c)%2 == 0 {
				checker[r][c] = '1'
			} else {
				checker[r][c] = '0'
			}
		}
	}

	cases := []struct {
		name string
		grid [][]byte
		want int
	}{
		{"one solid island", allLand, 1},
		{"checkerboard", checker, n * n / 2},
	}

	for _, c := range cases {
		start := time.Now()
		got := numIslands(c.grid)
		elapsed := time.Since(start)

		if got != c.want {
			t.Errorf("%s (%dx%d): numIslands = %d, want %d", c.name, n, n, got, c.want)
		}
		if elapsed > budget {
			t.Errorf("%s (%dx%d): took %v, over the %v budget", c.name, n, n, elapsed.Round(time.Millisecond), budget)
		}
		t.Logf("%s (%dx%d): %v", c.name, n, n, elapsed.Round(time.Millisecond))
	}
}
