package invertbinarytree

import (
	"slices"
	"testing"
)

// cloneTree copies structure and values, so a test can keep the original shape
// around while the real tree is inverted in place.
func cloneTree(node *TreeNode) *TreeNode {
	if node == nil {
		return nil
	}
	return &TreeNode{Val: node.Val, Left: cloneTree(node.Left), Right: cloneTree(node.Right)}
}

// isMirror reports whether a and b are mirror images of each other: same values
// down the middle, with left and right swapped at every node. It compares two
// trees rather than producing one, so it never reimplements the inversion.
func isMirror(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Val == b.Val && isMirror(a.Left, b.Right) && isMirror(a.Right, b.Left)
}

func TestInvertTreeExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want []any
	}{
		// samples from the statement
		{"neetcode sample", []any{1, 2, 3, 4, 5, 6, 7}, []any{1, 3, 2, 7, 6, 5, 4}},
		{"leetcode sample 1", []any{4, 2, 7, 1, 3, 6, 9}, []any{4, 7, 2, 9, 6, 3, 1}},
		{"leetcode sample 2", []any{2, 1, 3}, []any{2, 3, 1}},
		{"leetcode sample 3", []any{}, []any{}},

		// degenerate
		{"single node", []any{1}, []any{1}},
		{"left child only", []any{1, 2}, []any{1, nil, 2}},
		{"right child only", []any{1, nil, 2}, []any{1, 2}},

		// gaps at several levels
		{"gappy", []any{1, 2, 3, nil, 4}, []any{1, 3, 2, nil, nil, 4}},
		{"deeper gaps", []any{1, 2, 3, 4, nil, nil, 5}, []any{1, 3, 2, 5, nil, nil, 4}},
	}

	for _, c := range cases {
		got := invertTree(buildTree(c.vals))

		if !treesEqual(got, buildTree(c.want)) {
			t.Errorf("%s: invertTree(%v) did not give %v", c.name, c.vals, c.want)
		}
	}
}

// The result must be the mirror of what went in, checked at every node rather
// than only at the root.
func TestInvertTreeIsAMirror(t *testing.T) {
	inputs := [][]any{
		{1},
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, nil, 4, 5, nil},
		{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8},
		{1, 2, nil, 3, nil, 4},
		{3, 9, 20, nil, nil, 15, 7},
	}

	for _, vals := range inputs {
		original := buildTree(vals)
		reference := cloneTree(original)

		got := invertTree(original)

		if !isMirror(got, reference) {
			t.Errorf("invertTree(%v) is not the mirror of the original", vals)
		}
	}
}

// Inverting twice puts the tree back exactly as it was.
func TestInvertTreeTwiceIsIdentity(t *testing.T) {
	inputs := [][]any{
		{1, 2, 3, 4, 5, 6, 7},
		{4, 2, 7, 1, 3, 6, 9},
		{1, 2, 3, nil, 4, nil, 5, 6},
		{1, 2, nil, 3, nil, 4},
	}

	for _, vals := range inputs {
		root := buildTree(vals)
		reference := cloneTree(root)

		if !treesEqual(invertTree(invertTree(root)), reference) {
			t.Errorf("inverting %v twice did not give the original back", vals)
		}
	}
}

// A BST inverted is sorted the other way: its inorder walk comes out descending.
func TestInvertTreeOfBSTReversesInorder(t *testing.T) {
	cases := [][]int{
		{5, 3, 8, 1, 4, 7, 9},
		{50, 30, 70, 20, 40, 60, 80},
		{1, 2, 3, 4, 5},
		{0, -10, 10, -20, -5, 5, 20},
	}

	for _, vals := range cases {
		ascending := slices.Clone(vals)
		slices.Sort(ascending)

		descending := slices.Clone(ascending)
		slices.Reverse(descending)

		got := inorderTraversal(invertTree(buildBST(vals)))
		if !slices.Equal(got, descending) {
			t.Errorf("BST from %v inverted: inorder = %v, want %v", vals, got, descending)
		}
	}
}

// The tree is rewired in place: the same nodes come back, with the same values,
// and the returned root is the node that was passed in.
func TestInvertTreeReusesNodes(t *testing.T) {
	root := buildTree([]any{10, 5, 15, 3, 7, 13, 18})

	nodes := map[*TreeNode]bool{}
	var collect func(n *TreeNode)
	collect = func(n *TreeNode) {
		if n == nil {
			return
		}
		nodes[n] = true
		collect(n.Left)
		collect(n.Right)
	}
	collect(root)
	countBefore := len(nodes)

	got := invertTree(root)

	if got != root {
		t.Error("invertTree returned a different node than the root it was given")
	}

	seen := 0
	var walk func(n *TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		if !nodes[n] {
			t.Errorf("node holding %d is not one of the original nodes", n.Val)
		}
		seen++
		walk(n.Left)
		walk(n.Right)
	}
	walk(got)

	if seen != countBefore {
		t.Errorf("tree has %d nodes after inverting, want %d", seen, countBefore)
	}
}

// nil in, nil out.
func TestInvertTreeEmpty(t *testing.T) {
	if got := invertTree(nil); got != nil {
		t.Errorf("invertTree(nil) = %v, want nil", got)
	}
}

// Constraints allow 100 nodes. A one-sided chain flips to the other side, and
// the depth stays the same.
func TestInvertTreeChains(t *testing.T) {
	const n = 100

	vals := make([]int, n)
	for i := range vals {
		vals[i] = i
	}

	left := chain(vals, true)
	got := invertTree(left)

	if !treesEqual(got, chain(vals, false)) {
		t.Error("a chain of 100 to the left did not become a chain of 100 to the right")
	}
	if depth := bfsDepth(got); depth != n {
		t.Errorf("depth after inverting = %d, want %d", depth, n)
	}
	if count := countNodes(got); count != n {
		t.Errorf("node count after inverting = %d, want %d", count, n)
	}

	right := chain(vals, false)
	if !treesEqual(invertTree(right), chain(vals, true)) {
		t.Error("a chain of 100 to the right did not become a chain of 100 to the left")
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

// treesEqual compares structure and values, used to prove the traversal is
// read-only.
func treesEqual(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Val == b.Val && treesEqual(a.Left, b.Left) && treesEqual(a.Right, b.Right)
}

func buildBST(vals []int) *TreeNode {
	var root *TreeNode
	for _, v := range vals {
		root = bstInsert(root, v)
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

func countNodes(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return 1 + countNodes(node.Left) + countNodes(node.Right)
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

func inorderTraversal(root *TreeNode) []int {
	arr := make([]int, 0)

	inorder(&arr, root)

	return arr
}

func inorder(arr *[]int, root *TreeNode) {
	if root != nil {
		inorder(arr, root.Left)
		*arr = append(*arr, root.Val)
		inorder(arr, root.Right)
	}
}
