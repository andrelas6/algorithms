package tree

import (
	"fmt"
	"strings"
	"testing"
)

// serialize is an independent oracle: write the tree out as text, with a marker
// for every missing child, then compare the two strings. Two trees are the same
// exactly when their serialisations match. It walks one tree at a time, so it
// shares nothing with a side-by-side comparison.
func serialize(node *TreeNode) string {
	var b strings.Builder

	var walk func(n *TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			b.WriteString("#,")
			return
		}
		fmt.Fprintf(&b, "%d,", n.Val)
		walk(n.Left)
		walk(n.Right)
	}
	walk(node)

	return b.String()
}

func TestIsSameTreeExamples(t *testing.T) {
	cases := []struct {
		name string
		p, q []any
		want bool
	}{
		// samples from the statement
		{"identical", []any{1, 2, 3}, []any{1, 2, 3}, true},
		{"same values, different sides", []any{1, 2}, []any{1, nil, 2}, false},
		{"same values, swapped children", []any{1, 2, 1}, []any{1, 1, 2}, false},
		{"children swapped", []any{1, 2, 3}, []any{1, 3, 2}, false},
		{"one child on different sides", []any{4, 7}, []any{4, nil, 7}, false},

		// degenerate
		{"both empty", []any{}, []any{}, true},
		{"first empty", []any{}, []any{1}, false},
		{"second empty", []any{1}, []any{}, false},
		{"single node, same value", []any{1}, []any{1}, true},
		{"single node, different value", []any{1}, []any{2}, false},
	}

	for _, c := range cases {
		p, q := buildTree(c.p), buildTree(c.q)

		if got := isSameTree(p, q); got != c.want {
			t.Errorf("%s: isSameTree(%v, %v) = %v, want %v", c.name, c.p, c.q, got, c.want)
		}
	}
}

// Same shape, one value different. The difference is placed in a different spot
// each time, including the very last node visited.
func TestIsSameTreeValueDifferences(t *testing.T) {
	base := []any{1, 2, 3, 4, 5, 6, 7}

	for i := range base {
		changed := make([]any, len(base))
		copy(changed, base)
		changed[i] = 99

		if isSameTree(buildTree(base), buildTree(changed)) {
			t.Errorf("value at level-order position %d changed to 99: got true, want false", i)
		}
		if isSameTree(buildTree(changed), buildTree(base)) {
			t.Errorf("value at level-order position %d changed to 99 (arguments swapped): got true, want false", i)
		}
	}
}

// Same values, different shape. Value-only comparison, or comparing traversals
// without markers for missing children, passes several of these.
func TestIsSameTreeShapeDifferences(t *testing.T) {
	cases := []struct {
		name string
		p, q []any
	}{
		{"chain left against chain right", []any{1, 2, nil, 3}, []any{1, nil, 2, nil, 3}},
		{"one node deeper", []any{1, 2, 3}, []any{1, 2, 3, 4}},
		{"extra node at the bottom", []any{1, 2, 3, 4, 5}, []any{1, 2, 3, 4, 5, 6}},
		{"mirror image", []any{1, 2, 3, 4, nil, nil, 5}, []any{1, 3, 2, 5, nil, nil, 4}},
		{"prefix of the other", []any{1, 2}, []any{1, 2, 3}},
		{"same inorder, different shape", []any{2, 1, 3}, []any{1, nil, 2, nil, 3}},
	}

	for _, c := range cases {
		if isSameTree(buildTree(c.p), buildTree(c.q)) {
			t.Errorf("%s: isSameTree(%v, %v) = true, want false", c.name, c.p, c.q)
		}
		if isSameTree(buildTree(c.q), buildTree(c.p)) {
			t.Errorf("%s (arguments swapped): got true, want false", c.name)
		}
	}
}

// Cross-check against the serialisation oracle, both ways round.
func TestIsSameTreeAgainstSerialisation(t *testing.T) {
	trees := [][]any{
		{},
		{1},
		{1, 2},
		{1, nil, 2},
		{1, 2, 3},
		{1, 3, 2},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, nil, 4, 5, nil},
		{1, 2, nil, 3, nil, 4},
		{10, 5, 15, 3, 7, 13, 18},
		{0, 0, 0},
		{-1, -2, -3},
	}

	for i, pv := range trees {
		for j, qv := range trees {
			p, q := buildTree(pv), buildTree(qv)

			got := isSameTree(p, q)
			want := serialize(p) == serialize(q)

			if got != want {
				t.Errorf("trees %d and %d (%v vs %v): isSameTree = %v, serialisation says %v", i, j, pv, qv, got, want)
			}
		}
	}
}

// Values repeat inside these trees, so only position can tell them apart.
func TestIsSameTreeRepeatedValues(t *testing.T) {
	cases := []struct {
		name string
		p, q []any
		want bool
	}{
		{"all ones, same shape", []any{1, 1, 1}, []any{1, 1, 1}, true},
		{"all ones, different shape", []any{1, 1, 1}, []any{1, 1, nil, 1}, false},
		{"all zeroes, same shape", []any{0, 0, 0, 0}, []any{0, 0, 0, 0}, true},
		{"all zeroes, one missing", []any{0, 0, 0, 0}, []any{0, 0, 0}, false},
		{"negatives, identical", []any{-5, -10, -1}, []any{-5, -10, -1}, true},
		{"value limits", []any{10000, -10000}, []any{10000, -10000}, true},
		{"value limits, one off", []any{10000, -10000}, []any{10000, -9999}, false},
	}

	for _, c := range cases {
		if got := isSameTree(buildTree(c.p), buildTree(c.q)); got != c.want {
			t.Errorf("%s: isSameTree(%v, %v) = %v, want %v", c.name, c.p, c.q, got, c.want)
		}
	}
}

// A tree is always the same as itself and as a copy of itself.
func TestIsSameTreeReflexive(t *testing.T) {
	trees := [][]any{
		{1},
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, nil, 3, nil, 4},
		{10, 5, 15, 3, 7, nil, 18, 1, 4},
	}

	for _, vals := range trees {
		root := buildTree(vals)

		if !isSameTree(root, root) {
			t.Errorf("%v compared with itself: got false, want true", vals)
		}
		if !isSameTree(root, cloneTree(root)) {
			t.Errorf("%v compared with a copy: got false, want true", vals)
		}
	}
}

// Comparing is read-only, for both arguments.
func TestIsSameTreeDoesNotModifyTrees(t *testing.T) {
	pVals := []any{10, 5, 15, 3, 7, nil, 18}
	qVals := []any{10, 5, 15, 3, 7, nil, 19}

	p, q := buildTree(pVals), buildTree(qVals)
	pBefore, qBefore := serialize(p), serialize(q)

	isSameTree(p, q)

	if serialize(p) != pBefore {
		t.Error("isSameTree modified its first argument")
	}
	if serialize(q) != qBefore {
		t.Error("isSameTree modified its second argument")
	}
}

// Constraints allow 100 nodes. The difference is put at the very last node, so
// a comparison that stops early has to walk the whole tree to find it.
func TestIsSameTreeLarge(t *testing.T) {
	const n = 100

	vals := make([]any, n)
	for i := range vals {
		vals[i] = i + 1
	}

	if !isSameTree(buildTree(vals), buildTree(vals)) {
		t.Errorf("two identical complete trees of %d nodes: got false, want true", n)
	}

	changed := make([]any, n)
	copy(changed, vals)
	changed[n-1] = -1
	if isSameTree(buildTree(vals), buildTree(changed)) {
		t.Errorf("complete trees of %d nodes differing at the last node: got true, want false", n)
	}

	ints := make([]int, n)
	for i := range ints {
		ints[i] = i
	}
	if !isSameTree(chain(ints, true), chain(ints, true)) {
		t.Errorf("two identical chains of %d: got false, want true", n)
	}
	if isSameTree(chain(ints, true), chain(ints, false)) {
		t.Errorf("a chain of %d to the left against one to the right: got true, want false", n)
	}
}
