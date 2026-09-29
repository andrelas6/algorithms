package graph

/*
GREEN
Time: O(m*n)
Space: O(m*n) due to the map and grid
*/
func orangesRotting(grid [][]int) int {

	// find all the rotten fruits
	// dfs or bfs to visit every fresh fruit, mark visited, then continue until there aer no visited
	// fruits anymore
	rows, cols := len(grid), len(grid[0])

	visited := make(map[gridPosition]struct{})
	queue := make([]gridPosition, 0)
	countFresh := 0
	// find all rotten
	for r := range rows {
		for c := range cols {
			if grid[r][c] == 2 {
				queue = append(queue, gridPosition{r, c})
				visited[gridPosition{r, c}] = struct{}{}
			}
			if grid[r][c] == 1 {
				countFresh++
			}
		}
	}

	rottenCount := len(visited)

	minute := 0
	for len(queue) != 0 {
		length := len(queue)
		for range length {
			node := queue[0]
			queue = queue[1:]
			r := node[0]
			c := node[1]
			// for each direction
			visit(grid, r+1, c, visited, &queue)
			visit(grid, r-1, c, visited, &queue)
			visit(grid, r, c+1, visited, &queue)
			visit(grid, r, c-1, visited, &queue)
		}
		if len(queue) != 0 {
			minute++
		}
	}

	if countFresh != (len(visited) - rottenCount) {
		return -1
	} else {
		return minute
	}
}

func visit(grid [][]int, r, c int, visited map[gridPosition]struct{}, queue *[]gridPosition) bool {
	// is out of bounds
	if r >= len(grid) || r < 0 || c >= len(grid[0]) || c < 0 {
		return false
	}

	if grid[r][c] != 1 {
		return false
	}

	if _, ok := visited[gridPosition{r, c}]; ok {
		return false
	}

	*queue = append(*queue, gridPosition{r, c})
	visited[gridPosition{r, c}] = struct{}{}
	return true
}
