package heap

import "container/heap"

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type KthLargest struct {
	minHeap *MinHeap
	k       int
}

// RED
// Time: O(n * logn) - add is log n, but that's done for all elements
// Space: O(k)
// Link: https://neetcode.io/problems/kth-largest-integer-in-a-stream/question
func Constructor(k int, nums []int) KthLargest {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for _, n := range nums {
		heap.Push(minHeap, n)
	}
	for minHeap.Len() > k {
		heap.Pop(minHeap)
	}

	return KthLargest{minHeap: minHeap, k: k}
}

func (this *KthLargest) Add(val int) int {
	// log n
	heap.Push(this.minHeap, val)

	if this.minHeap.Len() > this.k {
		heap.Pop(this.minHeap)
	}

	return (*this.minHeap)[0]
}
