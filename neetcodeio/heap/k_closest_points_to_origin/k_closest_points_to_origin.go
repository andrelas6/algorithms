package kclosestpointstoorigin

import "container/heap"

type MinHeap [][]int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i][0] < h[j][0] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.([]int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

/*
Time: O(n log n) because of building the heap
Space: O(n)
*/
func kClosest(points [][]int, k int) [][]int {
	minHeap := &MinHeap{}
	// first pass, create minHeap
	for _, point := range points {
		distance := calcEuclidianDistance(point[0], point[1])
		*minHeap = append(*minHeap, []int{distance, point[0], point[1]})
	}
	heap.Init(minHeap)

	result := make([][]int, k)
	for i := range k {
		value := heap.Pop(minHeap).([]int)
		result[i] = []int{value[1], value[2]}
	}

	return result
}

func calcEuclidianDistance(x, y int) int {
	// taking square root makes the implementation a bit more complicated
	return x*x + y*y
}
