package findduplicate

func findDuplicate(nums []int) int {
	var slow, fast int

	// solution
	// use floyd's cycle detection to find intersection
	// use second loop with slow pointer to find where they meet
	// that's the duplicate
	// not a fun problem to solve in 30 min - it's more about memorizing the solution
	// O(n) time
	// O(1) space

	// find the intersection point
	for {
		slow = nums[slow]
		fast = nums[nums[fast]]
		if slow == fast {
			break
		}
	}

	var slow2 int

	// find where they meet - that's the duplicate
	for {
		slow = nums[slow]
		slow2 = nums[slow2]

		if slow == slow2 {
			break
		}
	}

	return slow
}
