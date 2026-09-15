package arrays

/*
 * YELLOW: I could easily implement the brute force solution with time O(n^2), but figuring out the improved approach was hard. I needed all the hints.
 * LINK: https://neetcode.io/problems/products-of-array-discluding-self/question
 * Time: O(n)
 * Space: O(1) because the returned value does not count as space
 */
func productExceptSelf(nums []int) []int {
	// brute force
	// loop n := range nums
	//    for each otherN, multiply
	// O(n2)
	// I have two loops - one to grab the element, the other to calculate
	// the product of each element except itself
	// optimize
	// how to remove the inner loop completely?
	// do one pass, calculate hte product of every element
	// then use division
	// but how to handle zero? use a hashmap with index -> zero element present
	// if there are more than 2 zeroes, everything is zero

	// 1st calculate prefix
	// starting prefix
	prefix := 1
	result := make([]int, len(nums))
	result[0] = 1
	for i := 1; i < len(nums); i++ {
		prevValue := nums[i-1]
		prefix *= prevValue
		result[i] = prefix
	}

	// starting postfix -> the last element is already set because postfix is 1
	postfix := 1
	// calculate postfixes and apply in-place (prefix * postfix)
	for i := len(result) - 2; i >= 0; i-- {
		// 2,4,6,8
		// i == 2, postfix = 8, result[2] = 8 * 8 = 64
		// 2,4,64,48
		// i == 1, postfix = 48, result[1] = 96
		// i == 0, postfix = 192, result[0] = 192 * 1 = 192
		postfix = nums[i+1] * postfix
		result[i] = postfix * result[i]
	}

	return result
}
