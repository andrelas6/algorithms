package graph

type visited map[gridPosition]struct{}

/*
 * YELLOW
 * Time: O(3 ^ (m * n)) visit every node one time at least once per path. It's 3^ and not 4^ because the previous node is not walked again.
 * Space: O(m * n) because of the hashmap
 * Link: https://neetcode.io/problems/matrixDFS/question
 */
func countPaths(grid [][]int) int {
	visited := make(visited)

	return dfsCountPaths(grid, 0, 0, visited)
}

func dfsCountPaths(grid [][]int, r, c int, v visited) int {
	// base case
	// 1 out of bounds
	// 2 the value I am looking for
	// 3 visited
	rows, cols := len(grid), len(grid[0])
	// out of bounds
	if (r < 0 || c < 0) || (r >= rows || c >= cols) {
		return 0
	}
	// not the target value
	if grid[r][c] == 1 {
		return 0
	}
	// visited
	if _, ok := v[gridPosition{r, c}]; ok {
		return 0
	}

	// calculate the value - reached final value
	if r == rows-1 && c == cols-1 {
		return 1
	}

	// add to visited position except the target - it must be visited again
	v[gridPosition{r, c}] = struct{}{}

	// dfs up, right, down, left
	up := dfsCountPaths(grid, r-1, c, v)
	right := dfsCountPaths(grid, r, c+1, v)
	down := dfsCountPaths(grid, r+1, c, v)
	left := dfsCountPaths(grid, r, c-1, v)

	// clean up from visited because recursion on this value is done and should support backtracking
	delete(v, gridPosition{r, c})

	return up + right + down + left

}
