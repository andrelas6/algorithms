package tree

import "testing"

// diameterByWalking is an independent oracle. It treats the tree as an
// undirected graph: build a neighbour list, then from every node walk outwards
// level by level to find the farthest node, and take the largest distance seen.
// No heights, no recursion, so a bug in the height/diameter bookkeeping cannot
// hide behind the same mistake here. O(n^2), fine for small trees.
func diameterByWalking(root *TreeNode) int {
	if root == nil {
		return 0
	}

	neighbours := map[*TreeNode][]*TreeNode{}
	var collect func(n *TreeNode)
	collect = func(n *TreeNode) {
		if n == nil {
			return
		}
		for _, child := range []*TreeNode{n.Left, n.Right} {
			if child != nil {
				neighbours[n] = append(neighbours[n], child)
				neighbours[child] = append(neighbours[child], n)
				collect(child)
			}
		}
	}
	collect(root)

	nodes := []*TreeNode{root}
	for n := range neighbours {
		if n != root {
			nodes = append(nodes, n)
		}
	}

	best := 0
	for _, start := range nodes {
		seen := map[*TreeNode]bool{start: true}
		frontier := []*TreeNode{start}

		for steps := 0; len(frontier) > 0; steps++ {
			if steps > best {
				best = steps
			}

			var next []*TreeNode
			for _, n := range frontier {
				for _, nb := range neighbours[n] {
					if !seen[nb] {
						seen[nb] = true
						next = append(next, nb)
					}
				}
			}
			frontier = next
		}
	}

	return best
}

func TestDiameterExamples(t *testing.T) {
	cases := []struct {
		name string
		vals []any
		want int
	}{
		// samples from the statement
		{"leetcode sample 1", []any{1, 2, 3, 4, 5}, 3},
		{"leetcode sample 2", []any{1, 2}, 1},
		{"neetcode sample", []any{1, nil, 2, 3, 4, 5}, 3},

		// degenerate
		{"single node", []any{1}, 0},
		{"empty", []any{}, 0},
		{"left child only", []any{1, 2}, 1},
		{"right child only", []any{1, nil, 2}, 1},

		// full trees: the longest path goes down one side and up the other
		{"full depth 2", []any{1, 2, 3}, 2},
		{"full depth 3", []any{1, 2, 3, 4, 5, 6, 7}, 4},
		{"full depth 4", []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 6},
	}

	for _, c := range cases {
		if got := diameterOfBinaryTree(buildTree(c.vals)); got != c.want {
			t.Errorf("%s: diameterOfBinaryTree(%v) = %d, want %d", c.name, c.vals, got, c.want)
		}
	}
}

// The trap: the longest path does not have to pass through the root. In every
// tree here the answer is strictly larger than height(left) + height(right) at
// the root, so a solution that only measures the root is wrong on all of them.
// The comparison against the root-only value is part of the test, so these stay
// honest.
func TestDiameterAwayFromTheRoot(t *testing.T) {
	// a node with two branches of equal depth, hanging under a root that has
	// nothing on its other side
	wideUnderOneChild := func() *TreeNode {
		return &TreeNode{Val: 1, Left: &TreeNode{
			Val:   2,
			Left:  chain([]int{3, 5}, true),
			Right: chain([]int{4, 6}, false),
		}}
	}

	// the same idea, mirrored onto the right of the root
	wideUnderOneChildMirrored := func() *TreeNode {
		return &TreeNode{Val: 1, Right: &TreeNode{
			Val:   2,
			Left:  chain([]int{3, 5}, true),
			Right: chain([]int{4, 6}, false),
		}}
	}

	// a deep, wide subtree on the left and a single node on the right
	heavyLeftLightRight := func() *TreeNode {
		return &TreeNode{
			Val: 0,
			Left: &TreeNode{
				Val:   100,
				Left:  chain([]int{1, 2, 3}, true),
				Right: chain([]int{4, 5, 6}, false),
			},
			Right: &TreeNode{Val: 7},
		}
	}

	cases := []struct {
		name string
		make func() *TreeNode
		want int
	}{
		{"wide subtree under the only child", wideUnderOneChild, 4},
		{"wide subtree under the only child, mirrored", wideUnderOneChildMirrored, 4},
		{"deep wide left, single node right", heavyLeftLightRight, 6},
	}

	for _, c := range cases {
		if got := diameterOfBinaryTree(c.make()); got != c.want {
			t.Errorf("%s: diameterOfBinaryTree = %d, want %d", c.name, got, c.want)
		}
		if oracle := diameterByWalking(c.make()); oracle != c.want {
			t.Errorf("%s: expected value %d disagrees with the walking oracle %d", c.name, c.want, oracle)
		}

		// the whole point: the root alone does not see this path
		root := c.make()
		throughRoot := bfsDepth(root.Left) + bfsDepth(root.Right)
		if throughRoot >= c.want {
			t.Errorf("%s: path through the root is %d and the answer is %d, so this case does not test anything",
				c.name, throughRoot, c.want)
		}
	}
}

// Cross-check against the walking oracle on a spread of shapes.
func TestDiameterAgainstWalking(t *testing.T) {
	trees := []*TreeNode{
		nil,
		buildTree([]any{1}),
		buildTree([]any{1, 2, 3}),
		buildTree([]any{1, 2, 3, 4, 5}),
		buildTree([]any{1, 2, 3, 4, 5, 6, 7}),
		buildTree([]any{1, 2, 3, nil, 4, 5, nil}),
		buildTree([]any{1, 2, nil, 3, nil, 4}),
		buildTree([]any{1, nil, 2, nil, 3, nil, 4}),
		buildTree([]any{10, 5, 15, 3, 7, 13, 18, 1, 4, 6, 8}),
		buildTree([]any{1, 2, 3, 4, nil, nil, 5, 6, nil, nil, 7}),
		buildTree([]any{3, 9, 20, nil, nil, 15, 7}),
		chain([]int{1, 2, 3, 4, 5}, true),
		chain([]int{1, 2, 3, 4, 5}, false),
		buildBST([]int{5, 3, 8, 1, 4, 7, 9}),
		buildBST([]int{50, 30, 70, 20, 40, 60, 80, 10, 25, 35}),
		buildBST([]int{1, 2, 3, 4, 5, 6, 7, 8}),
	}

	for i, root := range trees {
		got, want := diameterOfBinaryTree(root), diameterByWalking(cloneTree(root))
		if got != want {
			t.Errorf("tree %d: diameterOfBinaryTree = %d, walking oracle = %d", i, got, want)
		}
	}
}

// The diameter counts edges, so it is never more than one less than the number
// of nodes, and a tree of one node has diameter 0.
func TestDiameterBounds(t *testing.T) {
	trees := [][]any{
		{1},
		{1, 2},
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, nil, 3, nil, 4},
		{10, 5, 15, 3, 7, nil, 18},
	}

	for _, vals := range trees {
		root := buildTree(vals)
		diameter, nodes := diameterOfBinaryTree(root), countNodes(root)

		if diameter < 0 {
			t.Errorf("%v: diameter = %d, want at least 0", vals, diameter)
		}
		if diameter > nodes-1 {
			t.Errorf("%v: diameter = %d, more than %d edges available for %d nodes", vals, diameter, nodes-1, nodes)
		}
	}
}

// Measuring is read-only.
func TestDiameterDoesNotModifyTree(t *testing.T) {
	vals := []any{10, 5, 15, 3, 7, nil, 18, 1, 4}

	root := buildTree(vals)
	reference := buildTree(vals)

	diameterOfBinaryTree(root)

	if !treesEqual(root, reference) {
		t.Error("diameterOfBinaryTree modified the tree")
	}
}

// dfs is the recursive worker: it returns the HEIGHT in nodes while recording
// the widest span it has seen through the result pointer.
func TestDfsReturnsHeightAndRecordsDiameter(t *testing.T) {
	cases := []struct {
		name         string
		vals         []any
		wantHeight   int
		wantDiameter int
	}{
		{"nil", []any{}, 0, 0},
		{"single node", []any{1}, 1, 0},
		{"two levels", []any{1, 2}, 2, 1},
		{"full depth 3", []any{1, 2, 3, 4, 5, 6, 7}, 3, 4},
		{"chain of 4", []any{1, 2, nil, 3, nil, 4}, 4, 3},
	}

	for _, c := range cases {
		result := 0

		gotHeight := dfs(buildTree(c.vals), &result)
		if gotHeight != c.wantHeight {
			t.Errorf("%s: dfs returned height %d, want %d", c.name, gotHeight, c.wantHeight)
		}
		if result != c.wantDiameter {
			t.Errorf("%s: dfs recorded diameter %d, want %d", c.name, result, c.wantDiameter)
		}
	}
}

// Constraints allow 10^4 nodes. A chain has diameter n-1, and the recursion goes
// that deep.
func TestDiameterLong(t *testing.T) {
	const n = 10_000

	vals := make([]int, n)
	for i := range vals {
		vals[i] = i
	}

	if got := diameterOfBinaryTree(chain(vals, true)); got != n-1 {
		t.Errorf("chain of %d to the left: diameter = %d, want %d", n, got, n-1)
	}
	if got := diameterOfBinaryTree(chain(vals, false)); got != n-1 {
		t.Errorf("chain of %d to the right: diameter = %d, want %d", n, got, n-1)
	}

	// two chains of n/2 meeting at the root: the path runs down one and up the other
	root := &TreeNode{Val: 0}
	root.Left = chain(vals[:n/2], true)
	root.Right = chain(vals[:n/2], false)

	if got := diameterOfBinaryTree(root); got != n {
		t.Errorf("two chains of %d under one root: diameter = %d, want %d", n/2, got, n)
	}
}
