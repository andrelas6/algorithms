package binarytreerightsideviewer

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

// buildTree builds a tree from the level-order form used in the problem
// statement, where nil marks an absent child and children are only listed for
// nodes that exist. buildTree([]any{1,2,3,nil,4,5,nil}) is Example 2.
func buildTree(vals []any) *TreeNode {
	if len(vals) == 0 || vals[0] == nil {
		return nil
	}

	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}

	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]

		if i < len(vals) {
			if vals[i] != nil {
				node.Left = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, node.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != nil {
				node.Right = &TreeNode{Val: vals[i].(int)}
				queue = append(queue, node.Right)
			}
			i++
		}
	}

	return root
}

// chain builds a degenerate tree: a single path of len(vals) nodes, each hanging
// off its parent on the given side. This is the shape that turns the recursion
// depth into O(n) instead of O(log n).
func chain(vals []int, left bool) *TreeNode {
	if len(vals) == 0 {
		return nil
	}

	root := &TreeNode{Val: vals[0]}
	cur := root
	for _, v := range vals[1:] {
		next := &TreeNode{Val: v}
		if left {
			cur.Left = next
		} else {
			cur.Right = next
		}
		cur = next
	}

	return root
}

func buildBST(vals []int) *TreeNode {
	var root *TreeNode
	for _, v := range vals {
		root = bstInsert(root, v)
	}
	return root
}

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

// treesEqual compares structure and values, used to prove the traversal is
// read-only.
func treesEqual(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Val == b.Val && treesEqual(a.Left, b.Left) && treesEqual(a.Right, b.Right)
}

// bstInsert / buildBST give an independent oracle: the inorder traversal of a
// binary SEARCH tree is its values in ascending order, a fact that holds without
// reimplementing the traversal being tested.
func bstInsert(root *TreeNode, v int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: v}
	}
	if v < root.Val {
		root.Left = bstInsert(root.Left, v)
	} else {
		root.Right = bstInsert(root.Right, v)
	}
	return root
}
