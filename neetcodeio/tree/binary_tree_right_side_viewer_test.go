package tree

import (
	"slices"
	"testing"
)

// rightViewDFS is an independent oracle: a depth-first walk that visits the right
// child before the left, so the first node reached at each depth is the one seen
// from the right. Different algorithm from the level-by-level queue under test.
func rightViewDFS(root *TreeNode) []int {
	view := []int{}

	var walk func(node *TreeNode, depth int)
	walk = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(view) {
			view = append(view, node.Val)
		}
		walk(node.Right, depth+1)
		walk(node.Left, depth+1)
	}
	walk(root, 0)

	return view
}

func TestRightSideViewExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want []int
	}{
		// samples from the statement
		{"neetcode sample 1", []any{1, 2, 3}, []int{1, 3}},
		{"neetcode sample 2", []any{1, 2, 3, 4, 5, 6, 7}, []int{1, 3, 7}},
		{"leetcode sample 1", []any{1, 2, 3, nil, 5, nil, 4}, []int{1, 3, 4}},
		{"leetcode sample 2", []any{1, 2, 3, 4, nil, nil, nil, 5}, []int{1, 3, 4, 5}},
		{"leetcode sample 3", []any{1, nil, 3}, []int{1, 3}},

		// degenerate
		{"empty tree", []any{}, []int{}},
		{"single node", []any{1}, []int{1}},

		// one child only: whichever side exists is the one you see
		{"left child only", []any{1, 2}, []int{1, 2}},
		{"right child only", []any{1, nil, 2}, []int{1, 2}},
	}

	for _, c := range cases {
		if got := rightSideView(buildTree(c.vals)); !slices.Equal(got, c.want) {
			t.Errorf("%s: rightSideView(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// The trap in this problem: the visible node is the rightmost node of its LEVEL,
// not the right child of the node above. A solution that just follows .Right
// down from the root passes the full-tree samples and fails every one of these.
func TestRightSideViewLeftSubtreeShowsThrough(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want []int
	}{
		// 3 has no children, so 4 (under 2) is visible at depth 2
		{"left subtree deeper", []any{1, 2, 3, 4}, []int{1, 3, 4}},
		// 5 is 2's right child and still the rightmost at its level
		{"right child of the left subtree", []any{1, 2, 3, 4, 5}, []int{1, 3, 5}},
		// the right subtree ends early; the left carries on for two more levels
		{"left subtree much deeper", []any{1, 2, 3, 4, nil, nil, nil, 5, nil, 6}, []int{1, 3, 4, 5, 6}},
		// the view switches sides between levels: right, right, left, left
		{"switches sides", []any{1, 2, 3, nil, 4, 5, nil, 6, nil, nil, nil, nil, 7}, []int{1, 3, 5, 6, 7}},
		// zigzag: alternating left/right links, one node per level
		{"zigzag", []any{1, 2, nil, nil, 3, 4}, []int{1, 2, 3, 4}},
	}

	for _, c := range cases {
		if got := rightSideView(buildTree(c.vals)); !slices.Equal(got, c.want) {
			t.Errorf("%s: rightSideView(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// Values carry no meaning for visibility: negatives, zeroes and duplicates must
// come out exactly as stored.
func TestRightSideViewValues(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want []int
	}{
		{"negatives", []any{-1, -2, -3}, []int{-1, -3}},
		{"zeroes", []any{0, 0, 0, nil, 0}, []int{0, 0, 0}},
		{"duplicates", []any{7, 7, 7, 7, 7, 7, 7}, []int{7, 7, 7}},
		// constraints allow -100 <= Val <= 100
		{"range limits", []any{0, 100, -100, -100, 100}, []int{0, -100, 100}},
	}

	for _, c := range cases {
		if got := rightSideView(buildTree(c.vals)); !slices.Equal(got, c.want) {
			t.Errorf("%s: rightSideView(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// Degenerate chains: one node per level, so every node is visible, from either
// side.
func TestRightSideViewChains(t *testing.T) {
	for _, n := range []int{1, 2, 5, 100} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i + 1
		}

		if got := rightSideView(chain(vals, true)); !slices.Equal(got, vals) {
			t.Errorf("left chain of %d: got %v, want %v", n, got, vals)
		}
		if got := rightSideView(chain(vals, false)); !slices.Equal(got, vals) {
			t.Errorf("right chain of %d: got %v, want %v", n, got, vals)
		}
	}
}

// A complete tree filled in level order with 1..n: level k starts at 2^k, and its
// rightmost node is 2^(k+1)-1, or n on a last level that is only partly filled.
// Constraints allow 100 nodes; 1000 checks the queue well past that.
func TestRightSideViewCompleteTree(t *testing.T) {
	for _, n := range []int{1, 2, 3, 6, 7, 8, 100, 1000} {
		vals := make([]any, n)
		for i := range vals {
			vals[i] = i + 1
		}

		var want []int
		for last := 1; ; last = last*2 + 1 {
			if last >= n {
				want = append(want, n)
				break
			}
			want = append(want, last)
		}

		if got := rightSideView(buildTree(vals)); !slices.Equal(got, want) {
			t.Errorf("complete tree of %d: got %v, want %v", n, got, want)
		}
	}
}

// Cross-check against the right-first DFS on a spread of shapes, and check the
// view has exactly one value per level.
func TestRightSideViewAgainstDFS(t *testing.T) {
	roots := []*TreeNode{
		nil,
		buildTree([]any{1}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, nil, 4, 5, nil}),
		buildTree([]any{1, 2, 3, nil, 4, nil, 5, 6}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		buildTree([]any{3, 9, 20, nil, nil, 15, 7}),
		chain([]int{1, 2, 3, 4, 5}, true),
		chain([]int{1, 2, 3, 4, 5}, false),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
		buildBST([]int{50, 30, 70, 20, 40, 10, 45, 5}),
		buildBST([]int{42, 17, 93, 8, 25, 76, 100, -13, 61, 3}),
	}

	for i, root := range roots {
		got, want := rightSideView(root), rightViewDFS(root)
		if !slices.Equal(got, want) {
			t.Errorf("tree %d: rightSideView = %v, right-first DFS = %v", i, got, want)
		}
		if depth := bfsDepth(root); len(got) != depth {
			t.Errorf("tree %d: view has %d values, want one per level (%d)", i, len(got), depth)
		}
	}
}

// Taking the view is read-only.
func TestRightSideViewDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	rightSideView(root)

	if !treesEqual(root, reference) {
		t.Error("rightSideView modified the tree")
	}
}
