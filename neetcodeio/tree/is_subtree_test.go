package tree

import "testing"

// subtreeBySerialisation is an independent oracle: serialise the tree hanging
// off every node and compare those strings with the serialisation of subRoot.
// It reuses `serialize` from the same-tree tests and never compares two trees
// side by side, so a bug in the paired walk cannot hide behind it.
func subtreeBySerialisation(root, subRoot *TreeNode) bool {
	want := serialize(subRoot)

	var found bool
	var walk func(n *TreeNode)
	walk = func(n *TreeNode) {
		if n == nil || found {
			return
		}
		if serialize(n) == want {
			found = true
			return
		}
		walk(n.Left)
		walk(n.Right)
	}
	walk(root)

	return found
}

func TestIsSubtreeExamples(t *testing.T) {
	cases := []struct {
		name          string
		root, subRoot []any
		want          bool
	}{
		// samples from the statement
		{"leetcode sample 1", []any{3, 4, 5, 1, 2}, []any{4, 1, 2}, true},
		{"leetcode sample 2", []any{3, 4, 5, 1, 2, nil, nil, nil, nil, 0}, []any{4, 1, 2}, false},
		{"neetcode sample 1", []any{1, 2, 3, 4, 5}, []any{2, 4, 5}, true},
		{"neetcode sample 2", []any{1, 2, 3, 4, 5}, []any{2, 4}, false},

		// the whole tree counts as a subtree of itself
		{"same tree", []any{1, 2, 3}, []any{1, 2, 3}, true},
		{"single node against itself", []any{1}, []any{1}, true},

		// degenerate
		{"empty root", []any{}, []any{1}, false},
		{"single node, different value", []any{1}, []any{2}, false},
		{"subtree bigger than the tree", []any{1, 2}, []any{1, 2, 3}, false},
	}

	for _, c := range cases {
		if got := isSubtree(buildTree(c.root), buildTree(c.subRoot)); got != c.want {
			t.Errorf("%s: isSubtree(%v, %v) = %v, want %v", c.name, c.root, c.subRoot, got, c.want)
		}
	}
}

// The trap: a match has to take the node AND everything under it. Finding the
// right values with extra nodes hanging below is not a match.
func TestIsSubtreeMustTakeEverythingBelow(t *testing.T) {
	cases := []struct {
		name          string
		root, subRoot []any
	}{
		// 4 -> (1, 2) exists, but 1 has a child 0 underneath
		{"extra node under the match", []any{3, 4, 5, 1, 2, nil, nil, nil, nil, 0}, []any{4, 1, 2}},
		// the leaf 2 in root has children
		{"leaf in subtree has children in the tree", []any{1, 2, 3, 4, 5}, []any{1, 2, 3}},
		// shape matches down to one missing child
		{"one child too many", []any{1, 2, 3, nil, 4}, []any{1, 2, 3}},
		// values all present but spread differently
		{"same values, wrong shape", []any{1, 2, 3, 4}, []any{1, 2, 4}},
	}

	for _, c := range cases {
		root, sub := buildTree(c.root), buildTree(c.subRoot)

		if isSubtree(root, sub) {
			t.Errorf("%s: isSubtree(%v, %v) = true, want false", c.name, c.root, c.subRoot)
		}
		if subtreeBySerialisation(buildTree(c.root), buildTree(c.subRoot)) {
			t.Errorf("%s: the oracle disagrees, so this case does not test what it claims", c.name)
		}
	}
}

// Every subtree of a tree must be found, at every depth and on both sides.
func TestIsSubtreeFindsEveryNode(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, 13, 18, 1, 4}
	root := buildTree(vals)

	var nodes []*TreeNode
	var collect func(n *TreeNode)
	collect = func(n *TreeNode) {
		if n == nil {
			return
		}
		nodes = append(nodes, n)
		collect(n.Left)
		collect(n.Right)
	}
	collect(root)

	for _, node := range nodes {
		// a fresh copy, so the match cannot be by pointer identity
		if !isSubtree(root, cloneTree(node)) {
			t.Errorf("the subtree rooted at %d was not found", node.Val)
		}
	}
}

// Repeated values: the match may be the second or third place a value appears,
// and an early value match must not end the search.
func TestIsSubtreeRepeatedValues(t *testing.T) {
	cases := []struct {
		name          string
		root, subRoot []any
		want          bool
	}{
		// the 1 at the root has children; the deeper 1 is a leaf, which matches
		{"match is the deeper copy", []any{1, 1, nil, nil, 2}, []any{1, nil, 2}, true},
		{"leaf match among repeats", []any{1, 1, 1}, []any{1}, true},
		{"repeated shape, match on the right", []any{1, 2, 2, 3, nil, nil, 3}, []any{2, nil, 3}, true},
		{"repeated shape, no match", []any{1, 2, 2, 3, nil, nil, 3}, []any{2, 3, 3}, false},
		{"all zeroes, match", []any{0, 0, 0}, []any{0, 0, 0}, true},
		{"all zeroes, wrong size", []any{0, 0}, []any{0, 0, 0}, false},
	}

	for _, c := range cases {
		if got := isSubtree(buildTree(c.root), buildTree(c.subRoot)); got != c.want {
			t.Errorf("%s: isSubtree(%v, %v) = %v, want %v", c.name, c.root, c.subRoot, got, c.want)
		}
	}
}

// Cross-check against the serialisation oracle for every pair of a set of trees.
func TestIsSubtreeAgainstSerialisation(t *testing.T) {
	trees := [][]any{
		{1},
		{1, 2},
		{1, nil, 2},
		{1, 2, 3},
		{2, 4, 5},
		{3, 4, 5, 1, 2},
		{4, 1, 2},
		{1, 2, 3, 4, 5},
		{10, 5, 15, 3, 7, 13, 18},
		{5, 3, 7},
		{1, 1, 1},
		{0, 0, 0},
		{-1, -2, -3},
	}

	for i, rootVals := range trees {
		for j, subVals := range trees {
			got := isSubtree(buildTree(rootVals), buildTree(subVals))
			want := subtreeBySerialisation(buildTree(rootVals), buildTree(subVals))

			if got != want {
				t.Errorf("trees %d and %d (%v, %v): isSubtree = %v, serialisation says %v",
					i, j, rootVals, subVals, got, want)
			}
		}
	}
}

// Searching is read-only, for both trees.
func TestIsSubtreeDoesNotModifyTrees(t *testing.T) {
	rootVals := []any{10, 5, 15, 3, 7, nil, 18}
	subVals := []any{5, 3, 7}

	root, sub := buildTree(rootVals), buildTree(subVals)
	rootBefore, subBefore := serialize(root), serialize(sub)

	isSubtree(root, sub)

	if serialize(root) != rootBefore {
		t.Error("isSubtree modified the tree it searched")
	}
	if serialize(sub) != subBefore {
		t.Error("isSubtree modified the subtree it was looking for")
	}
}

// dfsIsSame compares two trees in step; it is the piece isSubtree leans on.
func TestDfsIsSame(t *testing.T) {
	cases := []struct {
		name string
		a, b []any
		want bool
	}{
		{"both empty", []any{}, []any{}, true},
		{"one empty", []any{1}, []any{}, false},
		{"identical", []any{1, 2, 3}, []any{1, 2, 3}, true},
		{"one value off", []any{1, 2, 3}, []any{1, 2, 4}, false},
		{"same values, mirrored", []any{1, 2, 3}, []any{1, 3, 2}, false},
		{"one node deeper", []any{1, 2}, []any{1, 2, nil, 3}, false},
	}

	for _, c := range cases {
		if got := dfsIsSame(buildTree(c.a), buildTree(c.b)); got != c.want {
			t.Errorf("%s: dfsIsSame(%v, %v) = %v, want %v", c.name, c.a, c.b, got, c.want)
		}
	}
}

// Constraints allow 100 nodes in the tree and 1000 in the tree being searched.
// A chain is the worst shape: every node has to be compared against the whole
// subtree before the search moves on.
func TestIsSubtreeLong(t *testing.T) {
	const n = 1000

	ints := make([]int, n)
	for i := range ints {
		ints[i] = i
	}
	root := chain(ints, true)

	// the last 50 values, as their own chain: this is exactly the bottom of root
	tail := chain(ints[n-50:], true)
	if !isSubtree(root, tail) {
		t.Error("the bottom 50 nodes of a chain of 1000 were not found")
	}

	// the same values, but hanging the other way
	if isSubtree(root, chain(ints[n-50:], false)) {
		t.Error("a chain going the other way was reported as a subtree")
	}

	// a chain of values that appear, but not consecutively at the bottom
	notTail := chain(ints[n-50:n-1], true)
	if isSubtree(root, notTail) {
		t.Error("a chain missing the final node was reported as a subtree")
	}

	// the whole chain against itself
	if !isSubtree(root, chain(ints, true)) {
		t.Error("a chain of 1000 was not found inside itself")
	}
}
