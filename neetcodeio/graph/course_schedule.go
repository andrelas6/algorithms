package graph

/*
 * RED - I didn't know how to do DFS in a graph problem. Also, I didn't
 * start with the intuition of drawing the graph with the edges. That made the whole difference for the solution.
 * Time: O(n + m) I GUESS, where n is the courses and m is the courses with no prerequisites
 * Space: O(n), where n is the amount of courses
 */
func canFinish(numCourses int, prerequisites [][]int) bool {
	// [a, b], a -> b
	// 0 -> 1 ALL GOOD
	// [0, 1], [1 , 0]
	// 0 -> 1
	// 1 -> 0
	// deadlock

	// assemble set to find out cycle
	visited := make(map[int]struct{})

	// assemble map to check the prerequisites for traversing
	courseToPrerequisite := make(map[int][]int)

	// fill courses then later add the prereqs
	for course := range numCourses {
		// 0 -> []
		// 1 -> []
		courseToPrerequisite[course] = []int{}
	}

	// assemble course -> prerequisites
	for _, comboCP := range prerequisites {
		// prerequisities is an awful name.. x IS NOT A prereq.
		// [0,1]
		course := comboCP[0]
		prereq := comboCP[1]
		val, _ := courseToPrerequisite[course]

		val = append(val, prereq)
		courseToPrerequisite[course] = val
		// 0 -> [1]
	}

	// for every course, chcek the path to null or cycle
	// that's dfs -> check every path
	for index := range numCourses {
		// 0
		// 1
		if !dfs(index, visited, courseToPrerequisite) {
			return false
		}
	}

	return true
}

// if no cycles, all good
func dfs(course int, visited map[int]struct{}, courseToPrerequisite map[int][]int) bool {
	// dont' add to visit if it has been already visited
	if _, ok := visited[course]; ok {
		return false
	}

	// grab all neighbors
	prereqs, _ := courseToPrerequisite[course]
	// 1 -> []

	if len(prereqs) == 0 {
		return true
	}

	// add to visited - only non terminating edges can be here
	// 0 -> {}
	visited[course] = struct{}{}

	// visit all neighbors
	// [1]
	for _, neighbor := range prereqs {
		if !dfs(neighbor, visited, courseToPrerequisite) {
			return false
		}
	}

	delete(visited, course)
	courseToPrerequisite[course] = []int{}
	return true
}
