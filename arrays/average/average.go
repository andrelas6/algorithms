package average

func countResponseTimeRegressions(responseTimes []int32) int32 {
	var result int32
	var sumPrevious int64
	// some numbers can extrapolate the limit of the int
	var quantityPrevious int64
	if len(responseTimes) != 0 {
		for index, respTime := range responseTimes {
			if index == 0 {
				continue
			} else {
				// sum all the previous
				// find how many previous elements
				// get the average
				// compare respTime with average
				// record if >
				// index 1
				sumPrevious += int64(responseTimes[index-1])
				quantityPrevious = int64(index)
				average := sumPrevious / quantityPrevious
				if int64(respTime) > int64(average) {
					result += 1
				}
			}
		}
	}

	return result
}
