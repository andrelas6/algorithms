package atlanticpacificwater

func pacificAtlantic(heights [][]int) [][]int {
	rows, cols := len(heights), len(heights[0])

	type GridPosition [2]int
	pac, atl := make(map[GridPosition]struct{}), make(map[GridPosition]struct{})

	var dfs func(r, c int, visit map[GridPosition]struct{}, prevHeight int)
	dfs = func(r, c int, visit map[GridPosition]struct{}, prevHeight int) {
		// check visit
		if _, ok := visit[GridPosition{r, c}]; ok {
			return
		}

		// is out of bounds
		if r < 0 || c < 0 || r >= rows || c >= cols {
			return
		}

		if heights[r][c] < prevHeight {
			return
		}
		visit[GridPosition{r, c}] = struct{}{}

		// dfs in all directions
		dfs(r+1, c, visit, heights[r][c])
		dfs(r-1, c, visit, heights[r][c])
		dfs(r, c+1, visit, heights[r][c])
		dfs(r, c-1, visit, heights[r][c])
	}

	// first and last row
	for c := range cols {
		dfs(0, c, pac, heights[0][c])
		dfs(rows-1, c, atl, heights[rows-1][c])
	}

	// first and last columns
	for r := range rows {
		dfs(r, 0, pac, heights[r][0])
		dfs(r, cols-1, atl, heights[r][cols-1])
	}

	results := make([][]int, 0)
	for i := range rows {
		for j := range cols {
			_, okP := pac[GridPosition{i, j}]
			_, okA := atl[GridPosition{i, j}]

			if okP && okA {
				results = append(results, []int{i, j})
			}
		}
	}
	return results
}
