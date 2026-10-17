package surroundedregions

func solve(board [][]byte) {
	// walk every cell

	// 1. DFS from border cells to set everything to a temp value
	// 2. change every O to X
	// 3. convert every temp value back to O

	rows, cols := len(board), len(board[0])

	var capture func(r, c int)
	capture = func(r, c int) {
		if r == rows || c == cols || r < 0 || c < 0 || board[r][c] != byte('O') {
			return
		}

		board[r][c] = byte('T')

		// dfs in all directions
		capture(r+1, c)
		capture(r-1, c)
		capture(r, c+1)
		capture(r, c-1)
	}

	// dfs on every border cell - take unsurrounded regions out
	for i := range rows {
		for j := range cols {
			// run dfs only on the border cells
			if i == 0 || j == 0 || i == rows-1 || j == cols-1 {
				capture(i, j)
			}
		}
	}

	// convert the rest to surrounded
	for i := range rows {
		for j := range cols {
			if board[i][j] == byte('O') {
				board[i][j] = byte('X')
			}
		}
	}

	// convert temp value back to O
	for i := range rows {
		for j := range cols {
			if board[i][j] == byte('T') {
				board[i][j] = byte('O')
			}
		}
	}
}
