package sorting

// Definition for a pair.
// type Pair struct {
//     Key   int
//     Value string
// }

// [1,2]
func mergeSort(pairs []Pair) []Pair {
	// var zero []Pair
	return mergeSortR(pairs)
	// return zero
}

func merge(left []Pair, right []Pair) []Pair {
	sortedSlice := make([]Pair, len(left)+len(right))

	// left [1]
	// right [2]
	var i, j, k int
	for j < len(right) && i < len(left) {
		if left[i].Key <= right[j].Key {
			sortedSlice[k] = left[i]
			i++
		} else {
			sortedSlice[k] = right[j]
			j++
		}
		k++
	}

	for i < len(left) {
		sortedSlice[k] = left[i]
		i++
		k++
	}

	for j < len(right) {
		sortedSlice[k] = right[j]
		j++
		k++
	}

	return sortedSlice
}

func mergeSortR(arr []Pair) []Pair {
	var result []Pair
	// len = 5 -> 2
	length := len(arr)
	if length <= 1 {
		return arr
	} else {
		// 3
		// 3 / 2 -> 1
		// 2 / 2 -> 1

		middle := length / 2
		left := mergeSortR(arr[:middle])
		right := mergeSortR(arr[middle:])
		result = merge(left, right)
	}

	return result
}
