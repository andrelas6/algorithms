package tree

import "testing"

// bfsDepth is an independent oracle: count levels iteratively with a queue.
// Different algorithm from the recursion under test, so a bug in maxDepth
// cannot hide behind the same mistake here.
func bfsDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	levels := 0
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levels++
		next := make([]*TreeNode, 0, len(queue)*2)
		for _, node := range queue {
			if node.Left != nil {
				next = append(next, node.Left)
			}
			if node.Right != nil {
				next = append(next, node.Right)
			}
		}
		queue = next
	}

	return levels
}

func TestMaxDepthExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want int
	}{
		// samples from the statement
		{"neetcode sample", []any{1, 2, 3, nil, nil, 4}, 3},
		{"leetcode sample", []any{3, 9, 20, nil, nil, 15, 7}, 3},
		{"root and one right child", []any{1, nil, 2}, 2},

		// degenerate
		{"empty tree", []any{}, 0},
		{"single node", []any{1}, 1},

		// one child only - each side gets to be the deeper one
		{"left child only", []any{1, 2}, 2},
		{"right child only", []any{1, nil, 2}, 2},

		// full trees
		{"full depth 2", []any{1, 2, 3}, 2},
		{"full depth 3", []any{1, 2, 3, 4, 5, 6, 7}, 3},

		// lopsided: the answer must come from the deeper side, not the first one
		{"deeper on the left", []any{1, 2, 3, 4, nil, nil, nil, 5}, 4},
		{"deeper on the right", []any{1, 2, 3, nil, nil, 4, 5, nil, nil, nil, 6}, 4},

		// gaps at several levels
		{"gappy", []any{1, 2, 3, nil, 4, nil, 5, 6}, 4},

		// negatives and duplicates are allowed by the constraints and must not
		// be confused with depth
		{"negative values", []any{-10, -20, -30, -40}, 3},
		{"duplicate values", []any{1, 1, 1, 1}, 3},
		{"zeroes", []any{0, 0, 0}, 2},
	}

	for _, c := range cases {
		if got := maxDepth(buildTree(c.vals)); got != c.want {
			t.Errorf("%s: maxDepth(%v) = %d, want %d", c.name, c.vals, got, c.want)
		}
	}
}

// Degenerate chains: depth equals node count, and recursion depth does too.
// A solution that counts levels only down one hardcoded side still passes the
// balanced cases and fails one of these.
func TestMaxDepthChains(t *testing.T) {
	for _, n := range []int{1, 2, 3, 10, 100} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i
		}

		if got := maxDepth(chain(vals, true)); got != n {
			t.Errorf("left chain of %d nodes: maxDepth = %d, want %d", n, got, n)
		}
		if got := maxDepth(chain(vals, false)); got != n {
			t.Errorf("right chain of %d nodes: maxDepth = %d, want %d", n, got, n)
		}
	}
}

// Cross-check against the iterative level count on a spread of shapes.
func TestMaxDepthAgainstBFS(t *testing.T) {
	roots := []*TreeNode{
		nil,
		buildTree([]any{1}),
		buildTree([]any{1, 2}),
		buildTree([]any{1, nil, 2}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, nil, 4, 5, nil}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		buildTree([]any{3, 9, 20, nil, nil, 15, 7}),
		chain([]int{1, 2, 3, 4, 5}, true),
		chain([]int{1, 2, 3, 4, 5}, false),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
		buildBST([]int{1, 2, 3, 4, 5, 6, 7, 8}),
		buildBST([]int{42, 17, 93, 8, 25, 76, 100, -13, 61, 3}),
	}

	for i, root := range roots {
		got, want := maxDepth(root), bfsDepth(root)
		if got != want {
			t.Errorf("tree %d: maxDepth = %d, bfsDepth = %d", i, got, want)
		}
	}
}

// The defining recurrence, checked at every node rather than only at the root:
// depth(node) == 1 + max(depth(left), depth(right)).
func TestMaxDepthRecurrenceHoldsEverywhere(t *testing.T) {
	roots := []*TreeNode{
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, nil, 4, nil, 5, 6}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4}),
		chain([]int{1, 2, 3, 4, 5}, true),
		buildBST([]int{50, 30, 70, 20, 40, 60, 80}),
	}

	var check func(t *testing.T, node *TreeNode)
	check = func(t *testing.T, node *TreeNode) {
		if node == nil {
			return
		}

		want := 1 + max(maxDepth(node.Left), maxDepth(node.Right))
		if got := maxDepth(node); got != want {
			t.Errorf("node %d: maxDepth = %d, want 1 + max(left, right) = %d", node.Val, got, want)
		}

		check(t, node.Left)
		check(t, node.Right)
	}

	for _, root := range roots {
		check(t, root)
	}
}

// Two bounds that must always hold: a non-empty tree is at least 1 deep, and no
// tree is deeper than it has nodes.
func TestMaxDepthBounds(t *testing.T) {
	roots := []*TreeNode{
		buildTree([]any{1}),
		buildTree([]any{1, 2, 3}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		chain([]int{1, 2, 3, 4, 5, 6}, false),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
	}

	for i, root := range roots {
		depth, nodes := maxDepth(root), countNodes(root)

		if depth < 1 {
			t.Errorf("tree %d: maxDepth = %d, want at least 1 for a non-empty tree", i, depth)
		}
		if depth > nodes {
			t.Errorf("tree %d: maxDepth = %d exceeds node count %d", i, depth, nodes)
		}
	}
}

// Computing the depth is read-only.
func TestMaxDepthDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	maxDepth(root)

	if !treesEqual(root, reference) {
		t.Error("maxDepth modified the tree")
	}
}

// depthFrom is the recursive worker and takes the level of its caller. Called
// with a nil node it reports that level unchanged, which is what makes the
// root's two calls start at 1.
func TestDepthFromStartingLevel(t *testing.T) {
	cases := []struct {
		name  string
		root  *TreeNode
		level int
		want  int
	}{
		{"nil at level 1", nil, 1, 1},
		{"nil at level 7", nil, 7, 7},
		{"leaf at level 1", buildTree([]any{1}), 1, 2},
		{"two levels from 1", buildTree([]any{1, 2}), 1, 3},
		{"offset start", buildTree([]any{1, 2}), 5, 7},
	}

	for _, c := range cases {
		if got := depthFrom(c.root, c.level); got != c.want {
			t.Errorf("%s: depthFrom(level=%d) = %d, want %d", c.name, c.level, got, c.want)
		}
	}
}
