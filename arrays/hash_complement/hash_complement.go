package hashcomplement

func findTaskPairForSlot(taskDurations []int32, slotLength int32) []int32 {
	seen := map[int32]int32{}

	for index, t := range taskDurations {
		// 1 4 2
		// 5
		// x 5 - 1 = 4, is 4 there? No, store 1 => 0
		// x 5 - 4 = 1, is 1 there? Yes, get 1 (0)
		x := slotLength - t

		v, ok := seen[x]

		if ok {
			return []int32{v, int32(index)}
		} else {
			seen[t] = int32(index)
		}
	}

	return []int32{-1, -1}

}
