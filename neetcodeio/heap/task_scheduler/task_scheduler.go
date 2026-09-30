package taskscheduler

import "container/heap"

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func leastInterval(tasks []byte, n int) int {
	// get count and add to max heap -> hash map
	// build slice, then init max heap with it
	// create queue for repetitive tasks
	// dance of max heap & queue

	// get count
	taskToCount := make(map[byte]int)

	for _, t := range tasks {
		taskToCount[t]++
	}

	// build slice and init max heap
	maxHeap := &MaxHeap{}
	for _, v := range taskToCount {
		*maxHeap = append(*maxHeap, v)
	}

	heap.Init(maxHeap)

	type TaskAndTime [2]int
	queue := make([]TaskAndTime, 0)

	// dance

	// if elements heap, need to process
	// if elements queue, need to process
	var timeTick int
	for maxHeap.Len() != 0 || len(queue) != 0 {
		timeTick++
		var taskCount int
		if maxHeap.Len() != 0 {
			taskCount = heap.Pop(maxHeap).(int)
			if taskCount > 1 {
				queue = append(queue, TaskAndTime{taskCount - 1, timeTick + n})
			}
		}

		// if there are tasks to be processed in the queue
		if len(queue) != 0 {
			head := queue[0]
			if head[1] == timeTick {
				heap.Push(maxHeap, head[0])
				// take tail of becuase it needs to be processed
				queue = queue[1:]
			}

		}

	}

	return timeTick
}
