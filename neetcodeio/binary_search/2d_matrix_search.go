package binarysearch

/*
 * YELLOW: I tried to calculate the midpoint of a 2d matrix, which was quite complicated. I needed the hint that searching the right row was the first right move.
 * then I had to do binary search on the row that contained the value, which I implemented correctly on the first go.
 * Time: O(log n) on the row search, then O(log m) on the binary search of the selected row
 * Space: O(1) because I don't use any data structures for this. I just use a pointer to find the selected row and the value.
 */
func searchMatrix(matrix [][]int, target int) bool {
	// binary search on the rows first
	// then another binary search on the selected row
	rows := len(matrix)
	cols := len(matrix[0])

	var row int = -1
	lastCol := cols - 1

	// replace with binary search
	left := 0
	right := rows - 1
	for left <= right {
		// get mid row
		// get left of mid and right of mid
		// compare
		// if target < left of mid, left stays, right is mid - 1
		// if target > right of mid, left is mid + 1, right stays
		midRow := (right + left) / 2
		leftOfMid := matrix[midRow][0]
		rightOfMid := matrix[midRow][lastCol]

		if target >= leftOfMid && target <= rightOfMid {
			row = midRow
			break
		}
		if target < leftOfMid {
			right = midRow - 1
		} else if target > rightOfMid {
			left = midRow + 1
		}
	}

	if row == -1 {
		return false
	}

	left = 0
	right = len(matrix[row]) - 1
	for left <= right {
		midPoint := (left + right) / 2
		midPointVal := matrix[row][midPoint]

		if target == midPointVal {
			return true
		}

		if target < midPointVal {
			right = midPoint - 1
		} else {
			left = midPoint + 1
		}
	}

	return false
}
