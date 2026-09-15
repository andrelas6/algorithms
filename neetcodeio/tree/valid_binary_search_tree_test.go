package tree

import (
	"math"
	"testing"
)

// inorderStrictlyIncreasing is an independent oracle: a binary tree is a valid
// BST exactly when its inorder traversal is strictly increasing. It reuses the
// already-tested inorderTraversal and never looks at bounds, so a bug in the
// boundary logic cannot hide behind the same mistake here.
func inorderStrictlyIncreasing(root *TreeNode) bool {
	vals := inorderTraversal(root)
	for i := 1; i < len(vals); i++ {
		if vals[i] <= vals[i-1] {
			return false
		}
	}
	return true
}

func TestIsValidBSTExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want bool
	}{
		// samples from the statement
		{"leetcode sample valid", []any{2, 1, 3}, true},
		{"leetcode sample invalid", []any{5, 1, 4, nil, nil, 3, 6}, false},
		{"neetcode sample valid", []any{2, 1, 3}, true},
		{"neetcode sample invalid", []any{1, 2, 3}, false},

		// degenerate
		{"empty tree", []any{}, true},
		{"single node", []any{1}, true},

		// one child only, on each side, both correct and wrong
		{"left child smaller", []any{2, 1}, true},
		{"left child larger", []any{2, 3}, false},
		{"right child larger", []any{2, nil, 3}, true},
		{"right child smaller", []any{2, nil, 1}, false},

		// the file's own diagram
		{"diagram tree", []any{4, 2, 5, 1, 3, nil, 6}, true},
		{"full depth 3", []any{4, 2, 6, 1, 3, 5, 7}, true},
		{"mirrored BST", []any{4, 6, 2, 7, 5, 3, 1}, false},
	}

	for _, c := range cases {
		if got := isValidBST(buildTree(c.vals)); got != c.want {
			t.Errorf("%s: isValidBST(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// The trap in this problem: every parent/child pair is locally fine, but a node
// deeper down breaks a bound set by an ancestor further up. A solution that only
// compares a node with its direct children passes every case here, when every
// one of them should fail.
func TestIsValidBSTAncestorBounds(t *testing.T) {
	cases := []struct {
		name string
		vals []any
	}{
		// 3 is a right grandchild of 5 but smaller than 5
		{"right subtree holds a smaller value", []any{5, 4, 6, nil, nil, 3, 7}},
		// 6 is a left grandchild of 5 but larger than 5
		{"left subtree holds a larger value", []any{5, 3, 7, nil, 6}},
		// 12 sits in 10's left subtree via 5 -> 8 -> 12
		{"violation three levels down", []any{10, 5, 15, 2, 8, nil, nil, nil, nil, 7, 12}},
		// the bound comes from the root, two turns away
		{"zigzag breaks root bound", []any{10, 5, nil, nil, 11}},
		{"zigzag breaks root bound on the right", []any{10, nil, 15, 9}},
	}

	for _, c := range cases {
		if isValidBST(buildTree(c.vals)) {
			t.Errorf("%s: isValidBST(%v) = true, want false (a descendant breaks an ancestor's bound)", c.name, c.vals)
		}
	}
}

// A BST in this problem is strict: equal values are allowed on neither side.
// This is what makes the comparison exclusive rather than inclusive.
func TestIsValidBSTDuplicatesAreInvalid(t *testing.T) {
	cases := []struct {
		name string
		vals []any
	}{
		{"equal left child", []any{1, 1}},
		{"equal right child", []any{1, nil, 1}},
		{"all equal", []any{2, 2, 2}},
		// equal to an ancestor, not the parent
		{"equals root deep on the left", []any{5, 3, nil, nil, 5}},
		{"equals root deep on the right", []any{5, nil, 8, 5}},
	}

	for _, c := range cases {
		if isValidBST(buildTree(c.vals)) {
			t.Errorf("%s: isValidBST(%v) = true, want false (duplicates are not allowed)", c.name, c.vals)
		}
	}
}

// Constraints allow -2^31 <= Val <= 2^31 - 1. Values at the limits must not be
// mistaken for the initial sentinel bounds.
func TestIsValidBSTValueLimits(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want bool
	}{
		{"single min int32", []any{math.MinInt32}, true},
		{"single max int32", []any{math.MaxInt32}, true},
		{"min and max as children", []any{0, math.MinInt32, math.MaxInt32}, true},
		{"max as left child", []any{0, math.MaxInt32}, false},
		{"min as right child", []any{0, nil, math.MinInt32}, false},
		{"duplicate max", []any{math.MaxInt32, nil, math.MaxInt32}, false},
		{"duplicate min", []any{math.MinInt32, math.MinInt32}, false},
		{"negatives", []any{-10, -20, -5, -30, -15}, true},
	}

	for _, c := range cases {
		if got := isValidBST(buildTree(c.vals)); got != c.want {
			t.Errorf("%s: isValidBST(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// Degenerate chains: the recursion depth equals the node count. Constraints
// allow up to 10^4 nodes.
func TestIsValidBSTChains(t *testing.T) {
	for _, n := range []int{2, 10, 10_000} {
		asc := make([]int, n)
		desc := make([]int, n)
		for i := range asc {
			asc[i] = i
			desc[i] = n - i
		}

		if !isValidBST(chain(asc, false)) {
			t.Errorf("ascending right chain of %d: got false, want true", n)
		}
		if !isValidBST(chain(desc, true)) {
			t.Errorf("descending left chain of %d: got false, want true", n)
		}
		if isValidBST(chain(asc, true)) {
			t.Errorf("ascending left chain of %d: got true, want false", n)
		}
		if isValidBST(chain(desc, false)) {
			t.Errorf("descending right chain of %d: got true, want false", n)
		}
	}
}

// Cross-check against the inorder oracle on a spread of valid and invalid
// shapes, including trees built by real BST insertion.
func TestIsValidBSTAgainstInorder(t *testing.T) {
	roots := []*TreeNode{
		nil,
		buildTree([]any{1}),
		buildTree([]any{2, 1, 3}),
		buildTree([]any{1, 2, 3}),
		buildTree([]any{5, 1, 4, nil, nil, 3, 6}),
		buildTree([]any{5, 4, 6, nil, nil, 3, 7}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 11}),
		buildTree([]any{3, 1, 5, 0, 2, 4, 6, nil, nil, nil, 3}),
		buildTree([]any{1, 1}),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
		buildBST([]int{1, 2, 3, 4, 5, 6, 7, 8}),
		buildBST([]int{42, 17, 93, 8, 25, 76, 100, -13, 61, 3}),
		chain([]int{5, 4, 3, 2, 1}, true),
		chain([]int{1, 2, 3, 2, 5}, false),
	}

	for i, root := range roots {
		got, want := isValidBST(root), inorderStrictlyIncreasing(root)
		if got != want {
			t.Errorf("tree %d (inorder %v): isValidBST = %v, inorder oracle = %v", i, inorderTraversal(root), got, want)
		}
	}
}

// Breaking a valid BST in one place, at any node, must be caught. Each node's
// value is swapped in turn for one just outside its allowed range.
func TestIsValidBSTCatchesEverySingleCorruption(t *testing.T) {
	vals := []int{50, 30, 70, 20, 40, 60, 80, 10, 25, 35, 45}

	var nodes []*TreeNode
	var collect func(node *TreeNode)
	collect = func(node *TreeNode) {
		if node == nil {
			return
		}
		nodes = append(nodes, node)
		collect(node.Left)
		collect(node.Right)
	}

	root := buildBST(vals)
	if !isValidBST(root) {
		t.Fatal("setup: uncorrupted BST reported invalid")
	}
	collect(root)

	for _, node := range nodes {
		original := node.Val
		for _, bad := range []int{math.MinInt32, math.MaxInt32} {
			node.Val = bad
			if isValidBST(root) == inorderStrictlyIncreasing(root) {
				node.Val = original
				continue
			}
			t.Errorf("node %d set to %d: isValidBST = %v, inorder oracle = %v",
				original, bad, isValidBST(root), inorderStrictlyIncreasing(root))
			node.Val = original
		}
	}
}

// Validation is read-only.
func TestIsValidBSTDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	isValidBST(root)

	if !treesEqual(root, reference) {
		t.Error("isValidBST modified the tree")
	}
}

// validWithDFS is the recursive worker. Its bounds are exclusive: a subtree is
// valid only if every value lies strictly between them.
func TestValidWithDFSBounds(t *testing.T) {
	cases := []struct {
		name        string
		vals        []any
		left, right int
		want        bool
	}{
		{"nil is always valid", []any{}, 5, 5, true},
		{"inside bounds", []any{5, 3, 7}, 0, 10, true},
		{"root equals left bound", []any{5}, 5, 10, false},
		{"root equals right bound", []any{5}, 0, 5, false},
		{"child equals left bound", []any{5, 1}, 1, 10, false},
		{"child equals right bound", []any{5, nil, 10}, 0, 10, false},
		{"child outside bounds", []any{5, nil, 11}, 0, 10, false},
	}

	for _, c := range cases {
		if got := validWithDFS(buildTree(c.vals), c.left, c.right); got != c.want {
			t.Errorf("%s: validWithDFS(%v, %d, %d) = %v, want %v", c.name, c.vals, c.left, c.right, got, c.want)
		}
	}
}

func TestCompareIsExclusive(t *testing.T) {
	cases := []struct {
		value, left, right int
		want               bool
	}{
		{5, 0, 10, true},
		{0, 0, 10, false},
		{10, 0, 10, false},
		{-1, 0, 10, false},
		{11, 0, 10, false},
		{5, 5, 5, false},
		{math.MinInt32, math.MinInt, math.MaxInt, true},
		{math.MaxInt32, math.MinInt, math.MaxInt, true},
	}

	for _, c := range cases {
		if got := compare(c.value, c.left, c.right); got != c.want {
			t.Errorf("compare(%d, %d, %d) = %v, want %v", c.value, c.left, c.right, got, c.want)
		}
	}
}
