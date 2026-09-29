package graph

func islandsAndTreasure(grid [][]int) {
	// dfs
	// path problem
	// recursion on each valid cell
	// when backtracking, set the amount of recursive calls
	// go forward
	// base case: -1, is visited (in recursion) and out of bounds

	rows, cols := len(grid), len(grid[0])
	visitedNodes := make(visited)
	queue := make([]gridPosition, 0)

	// do bfs from the treasures
	for r := range rows {
		for c := range cols {
			if grid[r][c] == 0 {
				queue = append(queue, gridPosition{r, c})
				visitedNodes[gridPosition{r, c}] = struct{}{}
			}
		}
	}

	bfsIslandTreasure(grid, visitedNodes, queue)
}

func bfsIslandTreasure(grid [][]int, visitedNodes visited, queue []gridPosition) {
	distance := 0
	for len(queue) != 0 {
		length := len(queue)
		for range length {
			position := queue[0]
			queue = queue[1:]

			r := position[0]
			c := position[1]
			grid[r][c] = distance

			if canAddToPath(grid, r+1, c, visitedNodes) {
				queue = append(queue, gridPosition{r + 1, c})
				visitedNodes[gridPosition{r + 1, c}] = struct{}{}
			}
			if canAddToPath(grid, r-1, c, visitedNodes) {
				queue = append(queue, gridPosition{r - 1, c})
				visitedNodes[gridPosition{r - 1, c}] = struct{}{}
			}
			if canAddToPath(grid, r, c+1, visitedNodes) {
				queue = append(queue, gridPosition{r, c + 1})
				visitedNodes[gridPosition{r, c + 1}] = struct{}{}
			}
			if canAddToPath(grid, r, c-1, visitedNodes) {
				queue = append(queue, gridPosition{r, c - 1})
				visitedNodes[gridPosition{r, c - 1}] = struct{}{}
			}
		}
		distance++
	}
}

func canAddToPath(grid [][]int, r, c int, visitedNodes visited) bool {
	if isOutOfBounds1(grid, r, c) ||
		isVisited1(r, c, visitedNodes) ||
		isWater(grid, r, c) {
		return false
	}

	return true
}

func isOutOfBounds1(grid [][]int, r, c int) bool {
	rows, cols := len(grid), len(grid[0])

	return r >= rows || r < 0 || c >= cols || c < 0
}
func isVisited1(r, c int, visitedNodes visited) bool {
	if _, ok := visitedNodes[gridPosition{r, c}]; ok {
		return true
	}

	return false
}
func isWater(grid [][]int, r, c int) bool {
	return grid[r][c] == -1
}
