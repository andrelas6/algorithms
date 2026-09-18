package tree

import "testing"

// balancedByHeights is an independent oracle: check every node using bfsDepth,
// the iterative level counter from the max-depth tests. It recomputes heights
// from scratch at each node (O(n^2)) and shares no code with the recursion under
// test.
func balancedByHeights(root *TreeNode) bool {
	if root == nil {
		return true
	}

	left, right := bfsDepth(root.Left), bfsDepth(root.Right)
	if left-right > 1 || right-left > 1 {
		return false
	}
	return balancedByHeights(root.Left) && balancedByHeights(root.Right)
}

func TestIsBalancedExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want bool
	}{
		// samples from the statement
		{"leetcode sample 1", []any{3, 9, 20, nil, nil, 15, 7}, true},
		{"leetcode sample 2", []any{1, 2, 2, 3, 3, nil, nil, 4, 4}, false},
		{"leetcode sample 3", []any{}, true},
		{"neetcode sample 1", []any{1, 2, 3, nil, nil, 4}, true},
		{"neetcode sample 2", []any{1, 2, 3, nil, nil, 4, nil, 5}, false},

		// degenerate
		{"single node", []any{1}, true},
		{"one child on the left", []any{1, 2}, true},
		{"one child on the right", []any{1, nil, 2}, true},
		{"two levels, full", []any{1, 2, 3}, true},
	}

	for _, c := range cases {
		if got := isBalanced(buildTree(c.vals)); got != c.want {
			t.Errorf("%s: isBalanced(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// A height difference of exactly 1 is allowed; 2 is not. These sit either side
// of that line.
func TestIsBalancedOnTheLimit(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want bool
	}{
		{"difference of 1 on the left", []any{1, 2, 3, 4}, true},
		{"difference of 1 on the right", []any{1, 2, 3, nil, nil, 4}, true},
		{"difference of 2 on the left", []any{1, 2, 3, 4, nil, nil, nil, 5}, false},
		{"difference of 2 on the right", []any{1, 2, 3, nil, nil, 4, nil, 5}, false},
		{"chain of three", []any{1, 2, nil, 3}, false},
		{"chain of three to the right", []any{1, nil, 2, nil, 3}, false},
	}

	for _, c := range cases {
		if got := isBalanced(buildTree(c.vals)); got != c.want {
			t.Errorf("%s: isBalanced(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// The trap: the root can look perfectly balanced while a node deep inside is
// not. Checking only the root's two heights passes all of these, and every one
// of them is unbalanced.
func TestIsBalancedDeepInside(t *testing.T) {
	// both sides of the root are 4 tall, so the root's difference is 0, but the
	// left side hides a node whose children differ by 2
	hiddenOnTheLeft := func() *TreeNode {
		return &TreeNode{
			Val:   1,
			Left:  &TreeNode{Val: 2, Left: chain([]int{3, 4, 5}, true)},
			Right: chain([]int{6, 7, 8}, false),
		}
	}

	hiddenOnTheRight := func() *TreeNode {
		return &TreeNode{
			Val:   1,
			Left:  chain([]int{6, 7, 8}, true),
			Right: &TreeNode{Val: 2, Right: chain([]int{3, 4, 5}, false)},
		}
	}

	// unbalanced two levels down, under an otherwise even tree
	hiddenTwoLevelsDown := func() *TreeNode {
		return &TreeNode{
			Val: 1,
			Left: &TreeNode{
				Val:   2,
				Left:  &TreeNode{Val: 4, Left: chain([]int{8, 9, 10}, true)},
				Right: &TreeNode{Val: 5},
			},
			Right: &TreeNode{
				Val:   3,
				Left:  chain([]int{6, 11, 12}, true),
				Right: chain([]int{7, 13, 14}, false),
			},
		}
	}

	cases := []struct {
		name string
		make func() *TreeNode
	}{
		{"hidden on the left", hiddenOnTheLeft},
		{"hidden on the right", hiddenOnTheRight},
		{"hidden two levels down", hiddenTwoLevelsDown},
	}

	for _, c := range cases {
		root := c.make()

		if isBalanced(root) {
			t.Errorf("%s: isBalanced = true, want false", c.name)
		}
		if balancedByHeights(c.make()) {
			t.Errorf("%s: the oracle disagrees, so this case does not test what it claims", c.name)
		}

		// the root itself looks fine, which is what makes this case worth having
		if left, right := bfsDepth(root.Left), bfsDepth(root.Right); left-right > 1 || right-left > 1 {
			t.Errorf("%s: the root already differs by more than 1 (%d vs %d), so the case is not testing depth",
				c.name, left, right)
		}
	}
}

// Cross-check against the oracle on a spread of shapes.
func TestIsBalancedAgainstHeights(t *testing.T) {
	trees := []*TreeNode{
		nil,
		buildTree([]any{1}),
		buildTree([]any{1, 2, 3}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}),
		buildTree([]any{1, 2, 3, nil, 4, 5, nil}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		buildTree([]any{3, 9, 20, nil, nil, 15, 7}),
		buildTree([]any{1, 2, 2, 3, 3, nil, nil, 4, 4}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		chain([]int{1, 2, 3}, true),
		chain([]int{1, 2, 3, 4, 5}, false),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
		buildBST([]int{1, 2, 3, 4, 5}),
		buildBST([]int{50, 30, 70, 20, 40, 60, 80, 10, 25, 35, 45}),
	}

	for i, root := range trees {
		got, want := isBalanced(root), balancedByHeights(cloneTree(root))
		if got != want {
			t.Errorf("tree %d: isBalanced = %v, height oracle = %v", i, got, want)
		}
	}
}

// Checking is read-only.
func TestIsBalancedDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	isBalanced(root)

	if !treesEqual(root, reference) {
		t.Error("isBalanced modified the tree")
	}
}

// dfsHeight returns the height in nodes and raises the flag as soon as any node
// below it is out of balance.
func TestDfsHeightAndFlag(t *testing.T) {
	cases := []struct {
		name          string
		vals          []any
		wantHeight    int
		wantUnbalance bool
	}{
		{"nil", []any{}, 0, false},
		{"single node", []any{1}, 1, false},
		{"two levels", []any{1, 2}, 2, false},
		{"full depth 3", []any{1, 2, 3, 4, 5, 6, 7}, 3, false},
		{"chain of three", []any{1, 2, nil, 3}, 3, true},
		{"unbalanced deeper down", []any{1, 2, 3, 4, nil, nil, nil, 5}, 4, true},
	}

	for _, c := range cases {
		var flag bool

		gotHeight := dfsHeight(buildTree(c.vals), &flag)
		if gotHeight != c.wantHeight {
			t.Errorf("%s: dfsHeight = %d, want %d", c.name, gotHeight, c.wantHeight)
		}
		if flag != c.wantUnbalance {
			t.Errorf("%s: unbalanced flag = %v, want %v", c.name, flag, c.wantUnbalance)
		}
	}
}

func TestCalcHeightDiff(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{0, 0, 0},
		{3, 3, 0},
		{4, 3, 1},
		{3, 4, 1},
		{7, 2, 5},
		{2, 7, 5},
	}

	for _, c := range cases {
		if got := calcHeightDiff(c.a, c.b); got != c.want {
			t.Errorf("calcHeightDiff(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

// Constraints allow 5000 nodes.
func TestIsBalancedLarge(t *testing.T) {
	const n = 5000

	// a complete tree filled in level order is balanced
	vals := make([]any, n)
	for i := range vals {
		vals[i] = i + 1
	}
	if !isBalanced(buildTree(vals)) {
		t.Errorf("complete tree of %d nodes: got false, want true", n)
	}

	// a chain is as unbalanced as it gets, and drives the recursion n deep
	ints := make([]int, n)
	for i := range ints {
		ints[i] = i
	}
	if isBalanced(chain(ints, true)) {
		t.Errorf("chain of %d to the left: got true, want false", n)
	}
	if isBalanced(chain(ints, false)) {
		t.Errorf("chain of %d to the right: got true, want false", n)
	}

	// balanced except for one long tail at the very bottom
	root := buildTree(vals)
	deepest := root
	for deepest.Left != nil {
		deepest = deepest.Left
	}
	deepest.Left = chain([]int{-1, -2, -3}, true)

	if isBalanced(root) {
		t.Error("complete tree with a chain hung off its deepest node: got true, want false")
	}
}
