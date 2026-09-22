package graph

type Node struct {
	Val       int
	Neighbors []*Node
}

/*
YELLOW - didn't know how to to dfs in an adjacency list
Time: O(e + v), where e is edge and v is vertice
Space: O(v) on the map, O(v) on the recursion?
Link: https://neetcode.io/problems/clone-graph/question
*/
func cloneGraph(node *Node) *Node {
	if node == nil {
		return nil
	}
	oldToNew := make(map[*Node]*Node)

	return dfsClone(node, oldToNew)
}

func dfsClone(node *Node, oldToNew map[*Node]*Node) *Node {
	// base case
	// node is already created
	val, ok := oldToNew[node]

	if ok {
		return val
	}

	copiedNode := Node{Val: node.Val, Neighbors: []*Node{}}
	oldToNew[node] = &copiedNode

	for _, neighbor := range node.Neighbors {
		copiedNode.Neighbors = append(copiedNode.Neighbors, dfsClone(neighbor, oldToNew))
	}

	return &copiedNode
}
