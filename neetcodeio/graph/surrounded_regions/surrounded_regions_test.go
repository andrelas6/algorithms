package surroundedregions

import (
	"slices"
	"testing"
)

func toBoard(rows []string) [][]byte {
	board := make([][]byte, len(rows))
	for i, r := range rows {
		board[i] = []byte(r)
	}
	return board
}

func toRows(board [][]byte) []string {
	rows := make([]string, len(board))
	for i, r := range board {
		rows[i] = string(r)
	}
	return rows
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name     string
		board    []string
		expected []string
	}{
		{
			"sample 1",
			[]string{"XXXX", "XOOX", "XXOX", "XOXX"},
			[]string{"XXXX", "XXXX", "XXXX", "XOXX"},
		},
		{"sample 2: single X", []string{"X"}, []string{"X"}},

		// degenerate
		{"single O on the border", []string{"O"}, []string{"O"}},
		{"single row", []string{"XOOXO"}, []string{"XOOXO"}},
		{"single column", []string{"X", "O", "O", "X"}, []string{"X", "O", "O", "X"}},
		{"2x2 all O", []string{"OO", "OO"}, []string{"OO", "OO"}},
		{"all X", []string{"XXX", "XXX", "XXX"}, []string{"XXX", "XXX", "XXX"}},

		// captured
		{"single O in the middle", []string{"XXX", "XOX", "XXX"}, []string{"XXX", "XXX", "XXX"}},
		{
			"two separate regions, both captured",
			[]string{"XXXXX", "XOXOX", "XXXXX"},
			[]string{"XXXXX", "XXXXX", "XXXXX"},
		},

		// not captured: connected to the border
		{
			"inner region joined to the border",
			[]string{"XXXX", "XOOX", "XOXX", "XOXX"},
			[]string{"XXXX", "XOOX", "XOXX", "XOXX"},
		},
		{
			"long snake reaching the border",
			[]string{"XXXXX", "XOOOX", "XXXOX", "XOOOX", "XOXXX"},
			[]string{"XXXXX", "XOOOX", "XXXOX", "XOOOX", "XOXXX"},
		},
		{"all O", []string{"OOO", "OOO", "OOO"}, []string{"OOO", "OOO", "OOO"}},

		// near-misses
		{
			"diagonal to the border doesn't count",
			[]string{"OXX", "XOX", "XXX"},
			[]string{"OXX", "XXX", "XXX"},
		},
		{
			"one region safe, one captured",
			[]string{"XXXXXX", "XOXXOX", "XXXXOX", "XXXXOX"},
			[]string{"XXXXXX", "XXXXOX", "XXXXOX", "XXXXOX"},
		},
		{
			"ring of O around an X, all on the border",
			[]string{"OOO", "OXO", "OOO"},
			[]string{"OOO", "OXO", "OOO"},
		},
		{
			"O inside a ring of X inside a ring of O",
			[]string{"OOOOO", "OXXXO", "OXOXO", "OXXXO", "OOOOO"},
			[]string{"OOOOO", "OXXXO", "OXXXO", "OXXXO", "OOOOO"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := toBoard(tt.board)
			solve(board)
			if got := toRows(board); !slices.Equal(got, tt.expected) {
				t.Errorf("solve(%v): expected %v, got %v", tt.board, tt.expected, got)
			}
		})
	}
}

// A 200x200 board of all O is one region of 40,000 cells, so the
// recursion goes very deep. Go grows its stack, so this should still pass.
func TestSolveLarge(t *testing.T) {
	const n = 200
	board := make([][]byte, n)
	for i := range board {
		board[i] = []byte(string(slices.Repeat([]byte{'O'}, n)))
	}
	solve(board)
	for i := range board {
		for j := range board[i] {
			if board[i][j] != 'O' {
				t.Fatalf("cell (%d,%d): expected O, got %c", i, j, board[i][j])
			}
		}
	}
}
