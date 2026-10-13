package copyrandomlist

import "testing"

// buildList builds a list from vals; randoms[i] is the index node i's Random
// points to, or -1 for nil.
func buildList(vals []int, randoms []int) *Node {
	nodes := make([]*Node, len(vals))
	for i, v := range vals {
		nodes[i] = &Node{Val: v}
	}
	for i := range nodes {
		if i+1 < len(nodes) {
			nodes[i].Next = nodes[i+1]
		}
		if randoms[i] >= 0 {
			nodes[i].Random = nodes[randoms[i]]
		}
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// flatten returns vals, random indices (-1 for nil), and the node pointers.
func flatten(head *Node) ([]int, []int, []*Node) {
	var nodes []*Node
	index := map[*Node]int{}
	for n := head; n != nil; n = n.Next {
		index[n] = len(nodes)
		nodes = append(nodes, n)
	}
	vals := make([]int, len(nodes))
	randoms := make([]int, len(nodes))
	for i, n := range nodes {
		vals[i] = n.Val
		randoms[i] = -1
		if n.Random != nil {
			r, ok := index[n.Random]
			if !ok {
				randoms[i] = -2 // points outside its own list
			} else {
				randoms[i] = r
			}
		}
	}
	return vals, randoms, nodes
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCopyRandomList(t *testing.T) {
	tests := []struct {
		name    string
		vals    []int
		randoms []int
	}{
		{"empty", []int{}, []int{}},
		{"single, nil random", []int{1}, []int{-1}},
		{"single, random to self", []int{1}, []int{0}},
		{"sample 1", []int{3, 7, 4, 5}, []int{-1, 3, 0, 1}},
		{"sample 2", []int{1, 2, 3, 4, 5}, []int{-1, 2, 4, -1, 1}},
		{"sample 3", []int{1, 2}, []int{1, 1}},
		{"all random nil", []int{1, 2, 3}, []int{-1, -1, -1}},
		{"all random to head", []int{1, 2, 3}, []int{0, 0, 0}},
		{"all random to tail", []int{1, 2, 3}, []int{2, 2, 2}},
		{"duplicate values, different randoms", []int{5, 5, 5}, []int{2, 0, 1}},
		{"negative values", []int{-100, 0, 100}, []int{1, -1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := buildList(tt.vals, tt.randoms)
			got := copyRandomList(head)

			gotVals, gotRandoms, gotNodes := flatten(got)
			_, _, origNodes := flatten(head)

			if !equalInts(gotVals, tt.vals) {
				t.Fatalf("vals: expected %v, got %v", tt.vals, gotVals)
			}
			if !equalInts(gotRandoms, tt.randoms) {
				t.Fatalf("randoms: expected %v, got %v (-2 = points into original list)", tt.randoms, gotRandoms)
			}

			orig := map[*Node]bool{}
			for _, n := range origNodes {
				orig[n] = true
			}
			for i, n := range gotNodes {
				if orig[n] {
					t.Fatalf("node %d is shared with the original, not a deep copy", i)
				}
			}

			// original must be untouched
			origVals, origRandoms, _ := flatten(head)
			if !equalInts(origVals, tt.vals) || !equalInts(origRandoms, tt.randoms) {
				t.Fatalf("original list was modified: vals %v randoms %v", origVals, origRandoms)
			}
		})
	}
}
