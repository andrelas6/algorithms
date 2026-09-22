package hashmaps

/*
 * Time: O(n)
 * Space: O(n)
 * Link: https://neetcode.io/problems/top-k-elements-in-list/question
 */
func topKFrequent(nums []int, k int) []int {
	// do one pass, count frequency
	// on each iteration, store result in a slice (tail has the max)
	// at the end, that should be ordered

	// do one pass, count frequency
	// 1 - 1, 2 - 2, 3 - 3
	// iterate on the map
	// run bucket sort for O(n) sorting
	// loop from end to begin to get the top K
	// return using slicing list[:k]
	// bucket sort must have size len(nums) + 1, because 0 is also value

	// iterate on 0..k, return values

	numToFrequency := make(map[int]int)

	for _, n := range nums {
		numToFrequency[n]++
	}

	sorted := bucketSort(numToFrequency, nums)

	result := make([]int, 0)
	var i int = len(sorted) - 1
	for len(result) <= k && i >= 0 {
		result = append(result, sorted[i]...)
		i--
	}

	return result[:k]
}

func bucketSort(numToFrequency map[int]int, nums []int) [][]int {
	size := len(nums)
	// 0 must be included
	sorted := make([][]int, size+1)

	for k, v := range numToFrequency {
		sorted[v] = append(sorted[v], k)
	}

	return sorted
}
