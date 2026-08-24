package heightpointer

import "testing"

func leaf(v int32) *Node             { return &Node{Val: v} }
func node(v int32, l, r *Node) *Node { return &Node{Val: v, Left: l, Right: r} }

func TestHeight(t *testing.T) {
	cases := []struct {
		name string
		tree *Node
		want int32
	}{
		{"empty tree", nil, 0},
		{"single node", leaf(10), 1},

		//   5
		//  /
		// 3
		{"root with one left child", node(5, leaf(3), nil), 2},

		//   5
		//    \
		//     8
		{"root with one right child", node(5, nil, leaf(8)), 2},

		//        4
		//      /   \
		//     2     6
		//    / \   / \
		//   1   3 5   7
		{
			"balanced, 7 nodes",
			node(4,
				node(2, leaf(1), leaf(3)),
				node(6, leaf(5), leaf(7)),
			),
			3,
		},

		// 4 -> 3 -> 2 -> 1, all left
		{
			"left chain of 4",
			node(4, node(3, node(2, leaf(1), nil), nil), nil),
			4,
		},

		// 1 -> 2 -> 3 -> 4, all right
		{
			"right chain of 4",
			node(1, nil, node(2, nil, node(3, nil, leaf(4)))),
			4,
		},

		//     5
		//    / \
		//   2   8
		//  /     \
		// 1       9
		{
			"deeper on the right",
			node(5,
				node(2, leaf(1), nil),
				node(8, nil, leaf(9)),
			),
			3,
		},

		//        8
		//      /   \
		//     4     9
		//    / \
		//   2   6
		//  /
		// 1
		{
			"lopsided left, depth 4",
			node(8,
				node(4, node(2, leaf(1), nil), leaf(6)),
				leaf(9),
			),
			4,
		},
	}

	for _, c := range cases {
		if got := height(c.tree); got != c.want {
			t.Errorf("%s: height() = %d, want %d", c.name, got, c.want)
		}
	}
}
