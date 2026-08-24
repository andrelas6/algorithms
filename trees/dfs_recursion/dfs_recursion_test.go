package dfsrecursion

import "testing"

func TestGetBinarySearchTreeHeight(t *testing.T) {
	cases := []struct {
		name       string
		values     []int32
		leftChild  []int32
		rightChild []int32
		want       int32
	}{
		// example from the statement
		{
			"balanced, 7 nodes",
			[]int32{4, 2, 6, 1, 3, 5, 7},
			[]int32{1, 3, 5, -1, -1, -1, -1},
			[]int32{2, 4, 6, -1, -1, -1, -1},
			3,
		},

		// samples
		{"single node", []int32{10}, []int32{-1}, []int32{-1}, 1},
		{"root with one left child", []int32{5, 3}, []int32{1, -1}, []int32{-1, -1}, 2},

		// degenerate
		{"empty tree", []int32{}, []int32{}, []int32{}, 0},
		{
			"left chain of 4",
			[]int32{4, 3, 2, 1},
			[]int32{1, 2, 3, -1},
			[]int32{-1, -1, -1, -1},
			4,
		},
		{
			"right chain of 4",
			[]int32{1, 2, 3, 4},
			[]int32{-1, -1, -1, -1},
			[]int32{1, 2, 3, -1},
			4,
		},

		// the deeper side is the RIGHT one
		{
			"deeper on the right",
			[]int32{5, 2, 8, 1, 9},
			[]int32{1, 3, -1, -1, -1},
			[]int32{2, -1, 4, -1, -1},
			3,
		},

		// child indices are NOT in walk order: following index+1 gives the wrong answer
		{
			"scrambled indices",
			[]int32{10, 15, 5, 7, 3},
			[]int32{2, -1, 4, -1, -1},
			[]int32{1, -1, 3, -1, -1},
			3,
		},

		// the deepest path is down the left, and only one branch is long
		{
			"lopsided left, depth 4",
			[]int32{8, 4, 9, 2, 6, 1},
			[]int32{1, 3, -1, 5, -1, -1},
			[]int32{2, 4, -1, -1, -1, -1},
			4,
		},
	}

	for _, c := range cases {
		if got := getBinarySearchTreeHeight(c.values, c.leftChild, c.rightChild); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}
