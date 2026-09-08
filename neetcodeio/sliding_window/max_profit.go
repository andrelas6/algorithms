package slidingwindow

func maxProfit(prices []int) int {
	var left, right int
	right = len(prices) - 1
	j := right

	var result int

	for left < right {
		sell := prices[j]
		buy := prices[left]
		profit := sell - buy
		if profit > 0 && profit > result {
			result = profit
		}
		j--
		if left == j {
			left++
			j = right
		}
	}

	return result
}
