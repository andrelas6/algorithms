package graph

type gridPosition [2]int

/*
 * RED - had no idea how to do bfs here and that bfs was the solution
 * Time: O(rows*columns)
 * Space: I have a visited map that will add all elements as they are visited
 * land, so maybe O(n * m)
 * problem: https://neetcode.io/problems/count-number-of-islands/question
 */
func numIslands(grid [][]byte) (islandCount int) {
	// iterate on 2d matrix until find "1"
	// add to count of islands
	// check neighbours

	rows := len(grid)
	cols := len(grid[0])
	visitedIslands := make(map[gridPosition]struct{})

	isIsland := func(rowToVisit, colToVisit int) bool {
		// out of bounds
		// is visited
		// is 1

		return (rowToVisit >= 0 && colToVisit >= 0) && (rowToVisit < rows && colToVisit < cols) && !isPositionVisited(rowToVisit, colToVisit, visitedIslands) && grid[rowToVisit][colToVisit] == byte('1')
	}

	up := [2]int{-1, 0}
	down := [2]int{1, 0}
	right := [2]int{0, 1}
	left := [2]int{0, -1}
	visitablePositions := [4]gridPosition{up, down, left, right}

	bfs := func(row, col int) {
		queue := make([]gridPosition, 0)
		queue = append(queue, gridPosition{row, col})
		visitedIslands[gridPosition{row, col}] = struct{}{}

		// only add to the queue if element is "1"
		for len(queue) != 0 {
			length := len(queue)

			for _ = range length {
				// grab head
				// add to visited
				head := queue[0]

				// dequeue first
				queue = queue[1:]
				// then check visited
				for _, visitable := range visitablePositions {
					rowPositionToVisit := visitable[0] + head[0]
					colPositionToVisit := visitable[1] + head[1]
					result := isIsland(rowPositionToVisit, colPositionToVisit)
					visitedIslands[gridPosition{rowPositionToVisit, colPositionToVisit}] = struct{}{}
					if result {
						land := gridPosition{rowPositionToVisit, colPositionToVisit}
						queue = append(queue, land)
					}
				}

			}
			// when to add to the queue?
			// neighbor must be '1'
			// neighbor cannot be out of bounds
			// neighbor cannot be an element that has been visited before
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			position := grid[i][j]
			if position == byte('1') && !isPositionVisited(i, j, visitedIslands) {
				islandCount++
				// do bfs to mark every visited position (they are part of the islands)
				bfs(i, j)
			}
		}
	}

	return islandCount
}

func isPositionVisited(rowToVisit, colToVisit int, visitedIslands map[gridPosition]struct{}) bool {
	if _, ok := visitedIslands[gridPosition{rowToVisit, colToVisit}]; ok {
		return true
	}

	return false
}
