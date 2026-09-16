package tree

import (
	"math/rand/v2"
	"testing"
)

// Both implementations answer the same question, so every table runs against
// both: the recursive one and the loop that replaced its stack frames.
var lcaImpls = []struct {
	name string
	fn   func(root, p, q *TreeNode) *TreeNode
}{
	{"lowestCommonAncestor", lowestCommonAncestor},
	{"lowestCommonAncestorMoreImproved", lowestCommonAncestorMoreImproved},
}

// nodeWithValue finds the node holding v, so tests can pass the real nodes from
// the tree rather than loose copies.
func nodeWithValue(root *TreeNode, v int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == v {
		return root
	}
	if found := nodeWithValue(root.Left, v); found != nil {
		return found
	}
	return nodeWithValue(root.Right, v)
}

// pathTo returns the nodes from root down to target, or nil if it isn't there.
func pathTo(root, target *TreeNode) []*TreeNode {
	if root == nil {
		return nil
	}
	if root == target {
		return []*TreeNode{root}
	}
	if below := pathTo(root.Left, target); below != nil {
		return append([]*TreeNode{root}, below...)
	}
	if below := pathTo(root.Right, target); below != nil {
		return append([]*TreeNode{root}, below...)
	}
	return nil
}

// lcaByPaths is an independent oracle: the lowest common ancestor is the last
// node the two root-to-node paths have in common. It uses no value comparison at
// all, so a bug in the left/right decisions cannot hide behind the same mistake.
func lcaByPaths(root, p, q *TreeNode) *TreeNode {
	pathP, pathQ := pathTo(root, p), pathTo(root, q)

	var lca *TreeNode
	for i := 0; i < len(pathP) && i < len(pathQ); i++ {
		if pathP[i] != pathQ[i] {
			break
		}
		lca = pathP[i]
	}
	return lca
}

// checkLCA runs the implementation on the nodes holding pVal and qVal and
// compares against want.
func checkLCA(t *testing.T, name string, root *TreeNode, pVal, qVal, want int) {
	t.Helper()

	p, q := nodeWithValue(root, pVal), nodeWithValue(root, qVal)
	if p == nil || q == nil {
		t.Fatalf("%s: setup is wrong, %d or %d is not in the tree", name, pVal, qVal)
	}

	for _, impl := range lcaImpls {
		got := impl.fn(root, p, q)
		if got == nil {
			t.Errorf("%s: %s(p=%d, q=%d) = nil, want %d", name, impl.name, pVal, qVal, want)
			continue
		}
		if got.Val != want {
			t.Errorf("%s: %s(p=%d, q=%d) = %d, want %d", name, impl.name, pVal, qVal, got.Val, want)
		}
	}
}

func TestLowestCommonAncestorExamples(t *testing.T) {
	// 5 / 3 8 / 1 4 7 9
	neetcode := buildTree([]any{5, 3, 8, 1, 4, 7, 9})
	// 6 / 2 8 / 0 4 7 9 / _ _ 3 5
	leetcode := buildTree([]any{6, 2, 8, 0, 4, 7, 9, nil, nil, 3, 5})

	cases := []struct {
		name             string
		root             *TreeNode
		pVal, qVal, want int
	}{
		// samples from the statement
		{"neetcode sample 1", neetcode, 3, 8, 5},
		{"neetcode sample 2", neetcode, 3, 4, 3},
		{"leetcode sample 1", leetcode, 2, 8, 6},
		{"leetcode sample 2", leetcode, 2, 4, 2},

		// the split happens at the root
		{"split at the root", neetcode, 1, 9, 5},
		{"split lower down", neetcode, 7, 9, 8},

		// p and q swapped: the answer must not depend on the order
		{"swapped", neetcode, 8, 3, 5},
		{"swapped, split lower", neetcode, 9, 7, 8},
	}

	for _, c := range cases {
		checkLCA(t, c.name, c.root, c.pVal, c.qVal, c.want)
	}
}

// The trap in this problem: a node is allowed to be its own ancestor. When one of
// the two nodes sits on the path to the other, that node IS the answer. A
// solution that always descends until the two split returns the wrong node here.
func TestLowestCommonAncestorOneIsAncestorOfTheOther(t *testing.T) {
	tree := buildTree([]any{6, 2, 8, 0, 4, 7, 9, nil, nil, 3, 5})

	cases := []struct {
		name             string
		pVal, qVal, want int
	}{
		{"root and a leaf", 6, 3, 6},
		{"root and its child", 6, 2, 6},
		{"parent and child", 2, 0, 2},
		{"grandparent and grandchild", 2, 3, 2},
		{"ancestor given second", 3, 2, 2},
		{"deep pair on one side", 4, 5, 4},
		{"same node twice", 4, 4, 4},
		{"same node twice at the root", 6, 6, 6},
	}

	for _, c := range cases {
		checkLCA(t, c.name, tree, c.pVal, c.qVal, c.want)
	}
}

// Degenerate shapes: a single node, and BSTs that are one long path.
func TestLowestCommonAncestorDegenerateShapes(t *testing.T) {
	single := buildTree([]any{1})
	checkLCA(t, "single node", single, 1, 1, 1)

	// inserted in ascending order: a chain going right
	rightChain := buildBST([]int{1, 2, 3, 4, 5})
	checkLCA(t, "right chain, ends", rightChain, 1, 5, 1)
	checkLCA(t, "right chain, middle pair", rightChain, 3, 5, 3)
	checkLCA(t, "right chain, same node", rightChain, 4, 4, 4)

	// descending: a chain going left
	leftChain := buildBST([]int{5, 4, 3, 2, 1})
	checkLCA(t, "left chain, ends", leftChain, 5, 1, 5)
	checkLCA(t, "left chain, middle pair", leftChain, 3, 1, 3)

	// two nodes only
	twoLeft := buildBST([]int{2, 1})
	checkLCA(t, "two nodes, left child", twoLeft, 1, 2, 2)

	twoRight := buildBST([]int{1, 2})
	checkLCA(t, "two nodes, right child", twoRight, 1, 2, 1)
}

// Negative values and zero are allowed by the constraints, and a run of values
// on one side of zero must not be confused with the other.
func TestLowestCommonAncestorNegativeValues(t *testing.T) {
	tree := buildBST([]int{0, -10, 10, -20, -5, 5, 20})

	cases := []struct {
		name             string
		pVal, qVal, want int
	}{
		{"both negative", -20, -5, -10},
		{"across zero", -20, 20, 0},
		{"zero with a negative", 0, -5, 0},
		{"both positive", 5, 20, 10},
	}

	for _, c := range cases {
		checkLCA(t, c.name, tree, c.pVal, c.qVal, c.want)
	}
}

// Every pair of nodes in a spread of BSTs, checked against the path oracle.
func TestLowestCommonAncestorEveryPairAgainstPaths(t *testing.T) {
	trees := [][]int{
		{5, 3, 8, 1, 4, 7, 9},
		{6, 2, 8, 0, 4, 7, 9, 3, 5},
		{50, 30, 70, 20, 40, 60, 80, 10, 25, 35, 45},
		{1, 2, 3, 4, 5, 6},          // right chain
		{6, 5, 4, 3, 2, 1},          // left chain
		{10, 5, 15, 3, 7, 13, 18, 1},
		{0, -10, 10, -20, -5, 5, 20},
	}

	for ti, vals := range trees {
		root := buildBST(vals)

		for _, pVal := range vals {
			for _, qVal := range vals {
				p, q := nodeWithValue(root, pVal), nodeWithValue(root, qVal)

				want := lcaByPaths(root, p, q)
				if want == nil {
					t.Fatalf("tree %d: setup is wrong, no path oracle answer for p=%d q=%d", ti, pVal, qVal)
				}

				for _, impl := range lcaImpls {
					got := impl.fn(root, p, q)
					if got == nil {
						t.Fatalf("tree %d: %s(p=%d, q=%d) = nil, path oracle = %d", ti, impl.name, pVal, qVal, want.Val)
					}
					if got.Val != want.Val {
						t.Errorf("tree %d: %s(p=%d, q=%d) = %d, path oracle = %d",
							ti, impl.name, pVal, qVal, got.Val, want.Val)
					}
				}
			}
		}
	}
}

// Random BSTs, every pair, against the path oracle. Fixed seed keeps it
// repeatable.
func TestLowestCommonAncestorRandomTrees(t *testing.T) {
	rng := rand.New(rand.NewPCG(13, 14))

	for range 200 {
		n := 1 + rng.IntN(40)

		vals := make([]int, 0, n)
		seen := map[int]bool{}
		for len(vals) < n {
			v := rng.IntN(200) - 100
			if !seen[v] {
				seen[v] = true
				vals = append(vals, v)
			}
		}

		root := buildBST(vals)

		for _, pVal := range vals {
			for _, qVal := range vals {
				p, q := nodeWithValue(root, pVal), nodeWithValue(root, qVal)

				want := lcaByPaths(root, p, q)

				for _, impl := range lcaImpls {
					got := impl.fn(root, p, q)
					if got == nil || got.Val != want.Val {
						t.Fatalf("vals=%v: %s(p=%d, q=%d) = %v, path oracle = %d",
							vals, impl.name, pVal, qVal, got, want.Val)
					}
				}
			}
		}
	}
}

// The answer must be a real node of the tree, not a copy, and it must actually
// have both p and q below it (or be one of them).
func TestLowestCommonAncestorReturnsANodeOfTheTree(t *testing.T) {
	root := buildBST([]int{50, 30, 70, 20, 40, 60, 80, 10, 25})
	vals := []int{50, 30, 70, 20, 40, 60, 80, 10, 25}

	inSubtree := func(node, target *TreeNode) bool {
		return pathTo(node, target) != nil
	}

	for _, pVal := range vals {
		for _, qVal := range vals {
			p, q := nodeWithValue(root, pVal), nodeWithValue(root, qVal)

			for _, impl := range lcaImpls {
				got := impl.fn(root, p, q)
				if got == nil {
					t.Fatalf("%s: p=%d q=%d: got nil", impl.name, pVal, qVal)
				}
				if pathTo(root, got) == nil {
					t.Errorf("%s: p=%d q=%d: returned a node with value %d that is not part of the tree", impl.name, pVal, qVal, got.Val)
					continue
				}
				if !inSubtree(got, p) || !inSubtree(got, q) {
					t.Errorf("%s: p=%d q=%d: returned %d, but it does not have both nodes below it", impl.name, pVal, qVal, got.Val)
				}
			}
		}
	}
}

// Finding an ancestor is read-only.
func TestLowestCommonAncestorDoesNotModifyTree(t *testing.T) {
	vals := []int{50, 30, 70, 20, 40, 60, 80}

	root := buildBST(vals)
	reference := buildBST(vals)

	for _, impl := range lcaImpls {
		impl.fn(root, nodeWithValue(root, 20), nodeWithValue(root, 40))

		if !treesEqual(root, reference) {
			t.Errorf("%s modified the tree", impl.name)
		}
	}
}

// Constraints allow up to 10^5 nodes. A chain that deep is where the recursion
// depth matters.
func TestLowestCommonAncestorDeepChain(t *testing.T) {
	const n = 100_000

	vals := make([]int, n)
	for i := range vals {
		vals[i] = i
	}
	// built directly rather than by repeated BST inserts, which would cost O(n^2)
	root := chain(vals, false) // ascending values down the right: still a valid BST

	checkLCA(t, "deep chain, ends", root, 0, n-1, 0)
	checkLCA(t, "deep chain, deepest pair", root, n-2, n-1, n-2)
	checkLCA(t, "deep chain, same node", root, n/2, n/2, n/2)
}
