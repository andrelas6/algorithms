package sorting

// Definition for a pair.
type Pair struct {
	Key   int
	Value string
}

func insertionSort(pairs []Pair) [][]Pair {
	// [k1, k2, k3, k3] -> sort based on Pair.key

	// two pointer technique - read value and then find its place in the array (and make it in place)
	// first loop -> get value
	// second loop -> iterate backwards until it finds hte place for "i"
	// base case -> when i reaches the end, leave first loop
	// base case -> when j > j - 1, leave second loop because place was found
	var result [][]Pair

	for i := range pairs {
		j := i
		// [4,3,2,1]
		// 4 j = 0 check, j - 1 > 0 not, skip
		// 3 j = 1, j - 1 = 0
		// 3 > 4 ? False, higher = 4, pairs[j-1] = 3, pairs[j] = 4, j = 0, skip [3, 4, 2, 1]
		// j = 2, 2 > 4? False, higher = 4, pair[j-1] = 2, pair[j] = 4, j = 1 [3,2,4,1]
		// j = 1, 2 > 3? False,
		for j >= 0 && j-1 >= 0 {
			// compare both
			// if j > j - 1, break
			// else, swap
			// higher = pairs[j - 1]
			// pairs[j - 1] = pairs[j]
			// pairs[j] = higher
			// decrease j j-=1
			previous := pairs[j-1]
			current := pairs[j]
			if current.Key >= previous.Key {
				break
			} else {
				pairs[j-1], pairs[j] = pairs[j], pairs[j-1]
				j -= 1
			}
		}

		// One insertion is done: element i has been walked into its slot in the
		// sorted region pairs[:i+1]. Record the array now -- once per element,
		// whether or not the inner loop moved anything. That is what makes the
		// state count exactly len(pairs).
		//
		// The copy is mandatory: appending `pairs` would store a slice header
		// pointing at the same backing array, so every later swap would rewrite
		// this "past" state too.
		snapshot := make([]Pair, len(pairs))
		copy(snapshot, pairs)
		result = append(result, snapshot)
	}

	return result
}
