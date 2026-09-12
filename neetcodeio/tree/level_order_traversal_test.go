package tree

import (
	"slices"
	"testing"
)

func levelsEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestLevelOrderExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want [][]int
	}{
		// samples from the statement
		{"neetcode sample", []any{1, 2, 3, nil, 4}, [][]int{{1}, {2, 3}, {4}}},
		{"leetcode sample", []any{3, 9, 20, nil, nil, 15, 7}, [][]int{{3}, {9, 20}, {15, 7}}},

		// degenerate
		{"empty tree", []any{}, nil},
		{"single node", []any{1}, [][]int{{1}}},

		// one child only - the level must still be a one-element slice
		{"left child only", []any{1, 2}, [][]int{{1}, {2}}},
		{"right child only", []any{1, nil, 2}, [][]int{{1}, {2}}},

		// full trees: every level complete
		{"full depth 2", []any{1, 2, 3}, [][]int{{1}, {2, 3}}},
		{"full depth 3", []any{1, 2, 3, 4, 5, 6, 7}, [][]int{{1}, {2, 3}, {4, 5, 6, 7}}},

		// left-to-right order within a level must be preserved
		{"ordered level", []any{1, 2, 3, 4, 5, 6, 7}, [][]int{{1}, {2, 3}, {4, 5, 6, 7}}},

		// gaps: children of different parents land on the same level, in order
		{"gappy", []any{1, 2, 3, nil, 4, 5, nil}, [][]int{{1}, {2, 3}, {4, 5}}},
		{"gappy deeper", []any{1, 2, 3, nil, 4, nil, 5, 6}, [][]int{{1}, {2, 3}, {4, 5}, {6}}},

		// values that must not be confused with structure
		{"negatives", []any{-1, -2, -3, -4}, [][]int{{-1}, {-2, -3}, {-4}}},
		{"duplicates", []any{1, 1, 1, 1}, [][]int{{1}, {1, 1}, {1}}},
		{"zeroes", []any{0, 0, 0}, [][]int{{0}, {0, 0}}},
	}

	for _, c := range cases {
		got := levelOrder(buildTree(c.vals))

		if !levelsEqual(got, c.want) {
			t.Errorf("%s: levelOrder(%v) = %v, want %v", c.name, c.vals, got, c.want)
		}
	}
}

// Degenerate chains: every level holds exactly one node, so the result is n
// levels of one value each. A solution that mixes up the level boundary
// collapses these into a single level.
func TestLevelOrderChains(t *testing.T) {
	for _, n := range []int{1, 2, 5, 100} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i
		}

		want := make([][]int, n)
		for i := range want {
			want[i] = []int{i}
		}

		if got := levelOrder(chain(vals, true)); !levelsEqual(got, want) {
			t.Errorf("left chain of %d: got %v, want %v", n, got, want)
		}
		if got := levelOrder(chain(vals, false)); !levelsEqual(got, want) {
			t.Errorf("right chain of %d: got %v, want %v", n, got, want)
		}
	}
}

// Structural properties, stated without reimplementing the traversal:
//   - one level per unit of depth
//   - every node appears exactly once
//   - each level is no more than twice the one above it
func TestLevelOrderStructure(t *testing.T) {
	roots := []*TreeNode{
		buildTree([]any{1}),
		buildTree([]any{1, 2, 3}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, nil, 4, nil, 5, 6}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		chain([]int{1, 2, 3, 4, 5}, true),
		chain([]int{1, 2, 3, 4, 5}, false),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
		buildBST([]int{42, 17, 93, 8, 25, 76, 100, -13, 61, 3}),
	}

	for i, root := range roots {
		got := levelOrder(root)

		if len(got) != maxDepth(root) {
			t.Errorf("tree %d: %d levels, want %d (the depth)", i, len(got), maxDepth(root))
		}

		total := 0
		for li, level := range got {
			if len(level) == 0 {
				t.Errorf("tree %d: level %d is empty; a level exists only if it has nodes", i, li)
			}
			if li > 0 && len(level) > 2*len(got[li-1]) {
				t.Errorf("tree %d: level %d has %d nodes, more than twice level %d's %d",
					i, li, len(level), li-1, len(got[li-1]))
			}
			total += len(level)
		}

		if n := countNodes(root); total != n {
			t.Errorf("tree %d: %d values across all levels, want %d (one per node)", i, total, n)
		}
	}
}

// Each level must hold exactly the nodes at that depth, left to right. This is
// checked by walking the tree independently and bucketing by depth.
func TestLevelOrderMatchesDepthBuckets(t *testing.T) {
	roots := []*TreeNode{
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{3, 9, 20, nil, nil, 15, 7}),
		buildTree([]any{1, 2, 3, nil, 4, nil, 5, 6}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4}),
		chain([]int{1, 2, 3, 4}, true),
		buildBST([]int{50, 30, 70, 20, 40, 60, 80}),
	}

	// collect values by depth, left to right, without using levelOrder
	var bucket func(node *TreeNode, d int, out *[][]int)
	bucket = func(node *TreeNode, d int, out *[][]int) {
		if node == nil {
			return
		}
		for len(*out) <= d {
			*out = append(*out, []int{})
		}
		(*out)[d] = append((*out)[d], node.Val)
		bucket(node.Left, d+1, out)
		bucket(node.Right, d+1, out)
	}

	for i, root := range roots {
		var want [][]int
		bucket(root, 0, &want)

		if got := levelOrder(root); !levelsEqual(got, want) {
			t.Errorf("tree %d: levelOrder = %v, want %v", i, got, want)
		}
	}
}

// The first level is always the root alone, and the last level is all leaves.
func TestLevelOrderFirstAndLastLevels(t *testing.T) {
	roots := []*TreeNode{
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{3, 9, 20, nil, nil, 15, 7}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
	}

	for i, root := range roots {
		got := levelOrder(root)

		if len(got) == 0 {
			t.Errorf("tree %d: no levels for a non-empty tree", i)
			continue
		}
		if !slices.Equal(got[0], []int{root.Val}) {
			t.Errorf("tree %d: first level = %v, want just the root %v", i, got[0], []int{root.Val})
		}

		leaves := map[int]bool{}
		var findLeaves func(n *TreeNode)
		findLeaves = func(n *TreeNode) {
			if n == nil {
				return
			}
			if n.Left == nil && n.Right == nil {
				leaves[n.Val] = true
			}
			findLeaves(n.Left)
			findLeaves(n.Right)
		}
		findLeaves(root)

		for _, v := range got[len(got)-1] {
			if !leaves[v] {
				t.Errorf("tree %d: %d is on the deepest level but is not a leaf", i, v)
			}
		}
	}
}

// Traversal is read-only.
func TestLevelOrderDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	levelOrder(root)

	if !treesEqual(root, reference) {
		t.Error("levelOrder modified the tree")
	}
}

// Constraints allow up to 2000 nodes. A perfect tree puts ~half of them on the
// last level, which is the widest the queue ever gets.
func TestLevelOrderWide(t *testing.T) {
	// perfect tree of depth 11 -> 2047 nodes, last level 1024
	depth := 11
	var build func(d, v int) *TreeNode
	build = func(d, v int) *TreeNode {
		if d == 0 {
			return nil
		}
		return &TreeNode{Val: v, Left: build(d-1, 2*v), Right: build(d-1, 2*v+1)}
	}

	root := build(depth, 1)
	got := levelOrder(root)

	if len(got) != depth {
		t.Fatalf("got %d levels, want %d", len(got), depth)
	}
	for i, level := range got {
		if want := 1 << i; len(level) != want {
			t.Errorf("level %d has %d nodes, want %d", i, len(level), want)
		}
	}
}
