package tree

import (
	"slices"
	"testing"
)

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

func buildBST(vals []int) *TreeNode {
	var root *TreeNode
	for _, v := range vals {
		root = bstInsert(root, v)
	}
	return root
}

// subtreeValues collects everything under node into a set. Collection order is
// irrelevant here, so this does not depend on inorder and can be used to check
// inorder's output.
func subtreeValues(node *TreeNode, out map[int]bool) {
	if node == nil {
		return
	}
	out[node.Val] = true
	subtreeValues(node.Left, out)
	subtreeValues(node.Right, out)
}

func countNodes(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return 1 + countNodes(node.Left) + countNodes(node.Right)
}

// treesEqual compares structure and values, used to prove the traversal is
// read-only.
func treesEqual(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Val == b.Val && treesEqual(a.Left, b.Left) && treesEqual(a.Right, b.Right)
}

func TestInorderTraversalExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want []int
	}{
		// full tree of depth 3
		{"example 1", []any{1, 2, 3, 4, 5, 6, 7}, []int{4, 2, 5, 1, 6, 3, 7}},
		// missing children on both sides
		{"example 2", []any{1, 2, 3, nil, 4, 5, nil}, []int{2, 4, 1, 5, 3}},
		// empty tree
		{"example 3", []any{}, []int{}},
	}

	for _, c := range cases {
		got := inorderTraversal(buildTree(c.vals))

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: inorderTraversal(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// Shapes chosen so left-visit, node-visit and right-visit each get to be the
// step that matters. A traversal with the three recursive steps in the wrong
// order (preorder or postorder by mistake) passes on a single node and fails on
// every one of these.
func TestInorderTraversalShapes(t *testing.T) {
	cases := []struct {
		name string
		root *TreeNode
		want []int
	}{
		{"nil root", nil, []int{}},
		{"single node", buildTree([]any{1}), []int{1}},

		// only a left child: the node is visited AFTER its child
		{"left child only", buildTree([]any{1, 2}), []int{2, 1}},
		// only a right child: the node is visited BEFORE its child
		{"right child only", buildTree([]any{1, nil, 2}), []int{1, 2}},

		{"both children", buildTree([]any{2, 1, 3}), []int{1, 2, 3}},

		// degenerate chains: recursion depth equals node count
		{"left chain", chain([]int{5, 4, 3, 2, 1}, true), []int{1, 2, 3, 4, 5}},
		{"right chain", chain([]int{1, 2, 3, 4, 5}, false), []int{1, 2, 3, 4, 5}},

		// zigzag: alternating left/right, so neither chain case covers it
		{"zigzag", buildTree([]any{1, 2, nil, nil, 3, 4}), []int{2, 4, 3, 1}},

		// left-heavy and right-heavy but not fully degenerate
		{"left heavy", buildTree([]any{4, 2, 5, 1, 3}), []int{1, 2, 3, 4, 5}},
		{"right heavy", buildTree([]any{2, 1, 4, nil, nil, 3, 5}), []int{1, 2, 3, 4, 5}},

		// negatives and zero: constraints allow -100 <= Val <= 100
		{"negatives", buildTree([]any{0, -50, 100, -100}), []int{-100, -50, 0, 100}},
		{"all negative", buildTree([]any{-1, -2, -3}), []int{-2, -1, -3}},

		// duplicate values are not excluded by the constraints
		{"duplicates", buildTree([]any{1, 1, 1}), []int{1, 1, 1}},

		// a deeper tree with gaps at several levels
		{"gappy", buildTree([]any{1, 2, 3, nil, 4, nil, 5, 6}), []int{2, 6, 4, 1, 3, 5}},
	}

	for _, c := range cases {
		got := inorderTraversal(c.root)

		if !slices.Equal(got, c.want) {
			t.Errorf("%s: inorderTraversal = %v, want %v", c.name, got, c.want)
		}
	}
}

// The empty tree must produce an empty slice, not nil -- the statement's expected
// output is [], and a nil slice serialises as null. The implementation's
// make([]int, 0) is what guarantees this; a bare `var arr []int` would pass
// slices.Equal (which treats nil and empty as equal) while changing the result.
func TestInorderTraversalEmptyIsNonNil(t *testing.T) {
	got := inorderTraversal(nil)

	if got == nil {
		t.Error("inorderTraversal(nil) returned a nil slice, want an empty non-nil slice so it renders as []")
	}
	if len(got) != 0 {
		t.Errorf("inorderTraversal(nil) = %v, want an empty slice", got)
	}
}

// Independent oracle: the inorder traversal of a binary search tree is its values
// in ascending order. Any mix-up in the visit order breaks this immediately, and
// the check does not depend on a second copy of the traversal logic.
func TestInorderTraversalOfBSTIsSorted(t *testing.T) {
	cases := [][]int{
		{5, 3, 8, 1, 4, 7, 9},
		{50, 30, 70, 20, 40, 60, 80},
		// inserted in ascending order -> degenerate right chain
		{1, 2, 3, 4, 5, 6, 7, 8},
		// descending -> degenerate left chain
		{8, 7, 6, 5, 4, 3, 2, 1},
		// negatives across zero
		{0, -10, 10, -20, -5, 5, 20},
		// random-ish
		{42, 17, 93, 8, 25, 76, 100, -13, 61, 3},
	}

	for _, vals := range cases {
		root := buildBST(vals)

		got := inorderTraversal(root)

		want := slices.Clone(vals)
		slices.Sort(want)

		if !slices.Equal(got, want) {
			t.Errorf("BST built from %v: inorder = %v, want %v (ascending)", vals, got, want)
		}
	}
}

// The defining property of inorder, stated without reimplementing it: for every
// node, every value in its LEFT subtree appears before it in the output, and
// every value in its RIGHT subtree appears after. Values are unique here so each
// one maps to a single output position.
func TestInorderTraversalSubtreeOrdering(t *testing.T) {
	roots := []*TreeNode{
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, nil, 4, 5, nil}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		chain([]int{1, 2, 3, 4, 5}, true),
		chain([]int{1, 2, 3, 4, 5}, false),
	}

	for ti, root := range roots {
		got := inorderTraversal(root)

		if n := countNodes(root); len(got) != n {
			t.Errorf("tree %d: traversal has %d values, want %d (one per node)", ti, len(got), n)
			continue
		}

		pos := make(map[int]int, len(got))
		for i, v := range got {
			pos[v] = i
		}

		var check func(node *TreeNode)
		check = func(node *TreeNode) {
			if node == nil {
				return
			}

			leftVals := map[int]bool{}
			subtreeValues(node.Left, leftVals)
			for v := range leftVals {
				if pos[v] > pos[node.Val] {
					t.Errorf("tree %d: %d is in the left subtree of %d but appears after it (positions %d vs %d) in %v",
						ti, v, node.Val, pos[v], pos[node.Val], got)
				}
			}

			rightVals := map[int]bool{}
			subtreeValues(node.Right, rightVals)
			for v := range rightVals {
				if pos[v] < pos[node.Val] {
					t.Errorf("tree %d: %d is in the right subtree of %d but appears before it (positions %d vs %d) in %v",
						ti, v, node.Val, pos[v], pos[node.Val], got)
				}
			}

			check(node.Left)
			check(node.Right)
		}
		check(root)
	}
}

// Traversal is read-only: it reads Val and follows Left/Right, and must not
// rewire or renumber anything.
func TestInorderTraversalDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	inorderTraversal(root)

	if !treesEqual(root, reference) {
		t.Error("inorderTraversal modified the tree")
	}
}

// `inorder` is the recursive worker: it APPENDS through a pointer to a slice
// rather than returning one. That contract matters -- it must add to whatever is
// already there instead of replacing it, which is what lets the two recursive
// calls and the append in between accumulate into one result.
func TestInorderAppendsToExistingSlice(t *testing.T) {
	arr := []int{99, 98}

	inorder(&arr, buildTree([]any{2, 1, 3}))

	want := []int{99, 98, 1, 2, 3}
	if !slices.Equal(arr, want) {
		t.Errorf("inorder appended to a non-empty slice: got %v, want %v", arr, want)
	}
}

// Constraints allow 100 nodes. Two shapes at that size: a balanced BST, and a
// fully degenerate chain where the recursion depth equals the node count.
func TestInorderTraversalLong(t *testing.T) {
	n := 100

	// balanced-ish BST via a level-order insert of a shuffled-but-deterministic
	// permutation; inorder must come out ascending
	vals := make([]int, n)
	for i := range vals {
		vals[i] = i - 50 // spans the allowed range, all distinct
	}
	shuffled := make([]int, 0, n)
	for i := 0; i < n; i += 7 {
		shuffled = append(shuffled, vals[i])
	}
	for i := range vals {
		if i%7 != 0 {
			shuffled = append(shuffled, vals[i])
		}
	}

	got := inorderTraversal(buildBST(shuffled))
	want := slices.Clone(vals)

	if !slices.Equal(got, want) {
		t.Errorf("BST of %d nodes: inorder is not ascending\ngot  %v\nwant %v", n, got, want)
	}

	// left chain of 100: recursion depth 100, output is the reverse of insertion
	desc := make([]int, n)
	for i := range desc {
		desc[i] = n - i
	}
	leftChain := chain(desc, true)

	got = inorderTraversal(leftChain)
	wantAsc := slices.Clone(desc)
	slices.Sort(wantAsc)

	if !slices.Equal(got, wantAsc) {
		t.Errorf("left chain of %d nodes: got %v, want %v", n, got, wantAsc)
	}

	// right chain of 100
	asc := make([]int, n)
	for i := range asc {
		asc[i] = i + 1
	}
	if got := inorderTraversal(chain(asc, false)); !slices.Equal(got, asc) {
		t.Errorf("right chain of %d nodes: got %v, want %v", n, got, asc)
	}
}
