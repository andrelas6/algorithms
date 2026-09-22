package graph

/*
GREEN
Time: O(m * n)
Space: O(m * n)
Link: https://neetcode.io/problems/max-area-of-island/question
*/
func maxAreaOfIsland(grid [][]int) int {
	// run dfs
	// if grid[i][j] is: 1) not 1 2) out of bounds and 3) visited, return 0
	// visited nodes remain saved, unlike finding "the path"
	visited := make(visitedPosition)
	result := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			newMax := dfsMaxArea(grid, i, j, visited)
			result = max(result, newMax)
		}
	}

	return result
}

type visitedPosition map[gridPosition]struct{}

func dfsMaxArea(grid [][]int, r, c int, visited visitedPosition) int {
	rows, cols := len(grid), len(grid[0])

	if isOutOfBounds(rows, cols, r, c) {
		return 0
	}
	if isVisited(r, c, visited) {
		return 0
	}

	visited[gridPosition{r, c}] = struct{}{}

	if isUnwanted(grid, r, c) {
		return 0
	}

	// go in all directions
	return 1 +
		dfsMaxArea(grid, r+1, c, visited) +
		dfsMaxArea(grid, r-1, c, visited) +
		dfsMaxArea(grid, r, c+1, visited) +
		dfsMaxArea(grid, r, c-1, visited)
}

func isOutOfBounds(maxRows, maxCols, r, c int) bool {
	return r >= maxRows || c >= maxCols || r < 0 || c < 0
}
func isVisited(r, c int, visited visitedPosition) bool {
	_, present := visited[gridPosition{r, c}]

	return present
}
func isUnwanted(grid [][]int, r, c int) bool {
	return grid[r][c] != 1
}
