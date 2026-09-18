package binarysearch

import (
	"slices"
	"testing"
)

func TestSearchMatrixExamples(t *testing.T) {
	neetcode := [][]int{
		{1, 2, 4, 8},
		{10, 11, 12, 13},
		{14, 20, 30, 40},
	}
	leetcode := [][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 60},
	}

	cases := []struct {
		name   string
		matrix [][]int
		target int
		want   bool
	}{
		// samples from the statement
		{"neetcode sample 1", neetcode, 10, true},
		{"neetcode sample 2", neetcode, 15, false},
		{"leetcode sample 1", leetcode, 3, true},
		{"leetcode sample 2", leetcode, 13, false},
	}

	for _, c := range cases {
		if got := searchMatrix(c.matrix, c.target); got != c.want {
			t.Errorf("%s: searchMatrix(target=%d) = %v, want %v", c.name, c.target, got, c.want)
		}
	}
}

// Every value in the matrix must be found, and the values just outside it must
// not be.
func TestSearchMatrixEveryValue(t *testing.T) {
	matrix := [][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 60},
	}

	for _, row := range matrix {
		for _, v := range row {
			if !searchMatrix(matrix, v) {
				t.Errorf("searchMatrix(target=%d) = false, want true (it is in the matrix)", v)
			}
			if searchMatrix(matrix, v-1) != slices.Contains(flatten(matrix), v-1) {
				t.Errorf("searchMatrix(target=%d) disagrees with whether %d is in the matrix", v-1, v-1)
			}
			if searchMatrix(matrix, v+1) != slices.Contains(flatten(matrix), v+1) {
				t.Errorf("searchMatrix(target=%d) disagrees with whether %d is in the matrix", v+1, v+1)
			}
		}
	}
}

func flatten(matrix [][]int) []int {
	out := []int{}
	for _, row := range matrix {
		out = append(out, row...)
	}
	return out
}

// The ends of each row, the ends of the matrix, and the gaps between rows.
func TestSearchMatrixBoundaries(t *testing.T) {
	matrix := [][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 60},
	}

	cases := []struct {
		name   string
		target int
		want   bool
	}{
		{"very first value", 1, true},
		{"very last value", 60, true},
		{"first value of the middle row", 10, true},
		{"last value of the middle row", 20, true},
		{"below everything", 0, false},
		{"far below everything", -10000, false},
		{"above everything", 61, false},
		{"far above everything", 10000, false},
		{"between two rows", 8, false},
		{"between two rows, just above a row end", 21, false},
		{"just below a row start", 9, false},
		{"inside a row but missing", 12, false},
	}

	for _, c := range cases {
		if got := searchMatrix(matrix, c.target); got != c.want {
			t.Errorf("%s: searchMatrix(target=%d) = %v, want %v", c.name, c.target, got, c.want)
		}
	}
}

func TestSearchMatrixSmallShapes(t *testing.T) {
	cases := []struct {
		name   string
		matrix [][]int
		target int
		want   bool
	}{
		{"single cell, present", [][]int{{5}}, 5, true},
		{"single cell, below", [][]int{{5}}, 4, false},
		{"single cell, above", [][]int{{5}}, 6, false},

		{"single row, first", [][]int{{1, 2, 3, 4}}, 1, true},
		{"single row, last", [][]int{{1, 2, 3, 4}}, 4, true},
		{"single row, middle", [][]int{{1, 2, 3, 4}}, 3, true},
		{"single row, missing", [][]int{{1, 3, 5, 7}}, 4, false},
		{"single row, above", [][]int{{1, 2, 3, 4}}, 5, false},

		{"single column, first", [][]int{{1}, {2}, {3}}, 1, true},
		{"single column, last", [][]int{{1}, {2}, {3}}, 3, true},
		{"single column, missing", [][]int{{1}, {3}, {5}}, 4, false},

		{"two by two, each value", [][]int{{1, 2}, {3, 4}}, 4, true},
		{"two by two, missing", [][]int{{1, 2}, {5, 6}}, 3, false},
	}

	for _, c := range cases {
		if got := searchMatrix(c.matrix, c.target); got != c.want {
			t.Errorf("%s: searchMatrix(target=%d) = %v, want %v", c.name, c.target, got, c.want)
		}
	}
}

// Constraints allow values from -10^4 to 10^4, so negatives and zero must work
// like any other value.
func TestSearchMatrixNegatives(t *testing.T) {
	matrix := [][]int{
		{-10000, -500, -3},
		{-2, -1, 0},
		{1, 99, 10000},
	}

	present := []int{-10000, -500, -3, -2, -1, 0, 1, 99, 10000}
	for _, v := range present {
		if !searchMatrix(matrix, v) {
			t.Errorf("searchMatrix(target=%d) = false, want true", v)
		}
	}

	for _, v := range []int{-9999, -501, -4, 2, 98, 100, 9999} {
		if searchMatrix(matrix, v) {
			t.Errorf("searchMatrix(target=%d) = true, want false", v)
		}
	}
}

// Searching must not rewrite the caller's matrix.
func TestSearchMatrixDoesNotModifyInput(t *testing.T) {
	matrix := [][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 60},
	}
	before := flatten(matrix)

	searchMatrix(matrix, 16)
	searchMatrix(matrix, 17)

	if after := flatten(matrix); !slices.Equal(after, before) {
		t.Errorf("searchMatrix modified the matrix: got %v, want %v", after, before)
	}
}

// Constraints allow 100 x 100. Values are 2 apart, so every odd value between
// them is a miss.
func TestSearchMatrixMaxSize(t *testing.T) {
	const m, n = 100, 100

	matrix := make([][]int, m)
	next := -10000
	for r := range m {
		matrix[r] = make([]int, n)
		for c := range n {
			matrix[r][c] = next
			next += 2
		}
	}

	for _, row := range matrix {
		for _, v := range row {
			if !searchMatrix(matrix, v) {
				t.Fatalf("searchMatrix(target=%d) = false, want true", v)
			}
			if searchMatrix(matrix, v+1) {
				t.Fatalf("searchMatrix(target=%d) = true, want false", v+1)
			}
		}
	}

	if searchMatrix(matrix, matrix[0][0]-1) {
		t.Error("value below the whole matrix reported as present")
	}
	if searchMatrix(matrix, matrix[m-1][n-1]+1) {
		t.Error("value above the whole matrix reported as present")
	}
}
