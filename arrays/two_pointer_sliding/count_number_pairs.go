package twopointersliding

func countAffordablePairs(prices []int32, budget int32) int32 {
	// Write your code here
	// [1,2,3,4,5] -> (1,2), (1,3), (1,4), (1,5), (2,3), (2,4), (2,5), (3,4)
	//
	// brute force
	// loop over every value -> then inner loop -> then count and store in another value O(n2)
	// if length = 0 or 1 -> early return
	// [1,2,3] 7
	// i = 1 -> loop j = 2 and j = 3 -> i <= j ? storeValue : skip
	// return withinBudgetCount
	//
	// optimise
	// two pointers that move together
	// [1,2,3,4] 7 (1,2) 1,3, 1,4 2,3 2,4 3,4
	// p1 = 1, p2 = 4
	// I know it's sorted -> so I find the two rightmost elements that p1 + p2 <= bugdget, then derive rest from length and index
	// p1 = 3, p2 = 4 -> p1 + p4 <= 7 ? yes, then return length - index + 1
	length := len(prices)

	var p1, p2 int32

	p2 = int32(length) - 1
	p1 = p2 - 1

	// 1,2,3,4,5
	// p1 = 3
	// p2 = 4
	//
	leftWindow := p1
	rightWindow := p2
	// 1, 2, 3.    3
	// p1 1 p2 2 -> 2 + 3 <= 3? No, leftWindow = 0, p1 = 0, p2 = 2
	// p1 0, p2 2 -> 1 + 3 <= 3? No, rightWindow = 1, leftWindow = -1
	count := 0
	for p1 >= 0 && p2 >= 1 {
		if prices[p1]+prices[p2] <= budget {
			// fmt.Printf("P1 = %d P2 = %d\n", p1,p2)
			// fmt.Printf("P1 = %d P2 = %d\n", p1,p2)
			// fmt.Printf("P1 = %d P2 = %d\n", p1,p2)
			count++
		}

		if p1 == 0 {
			rightWindow--
			leftWindow = rightWindow - 1
			p2 = rightWindow
			p1 = leftWindow
		} else {
			leftWindow--
			p1 = leftWindow
		}
	}

	return int32(count)
}
