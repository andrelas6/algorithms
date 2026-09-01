package sorting

// Definition for a pair.
// type Pair struct {
//     Key   int
//     Value string
// }

func QuickSortHelper(pairs []Pair, start int, end int) []Pair {
	if end-start <= 0 {
		return pairs
	}

	pivot := end
	left := start
	i := start

	pivot = end // adding to last element
	// [5,4,3,2] 2
	// [2,4,3,5]
	// [2,7,1,3]
	//  |
	//.   | [2, 1, 3, 7]
	for i != pivot {
		if pairs[i].Key < pairs[pivot].Key {
			pairs[left], pairs[i] = pairs[i], pairs[left]
			left++
		}
		i++
	}

	pairs[left], pairs[pivot] = pairs[pivot], pairs[left]

	QuickSortHelper(pairs, start, left-1)
	QuickSortHelper(pairs, left+1, end)

	return pairs
}

func QuickSort(pairs []Pair) []Pair {

	QuickSortHelper(pairs, 0, len(pairs)-1)

	return pairs
}
