package queue

import (
	"container/heap"
)

type PriorityQueue[T any] struct {
	items []*Item[T]
	less  func(a, b T) bool
}

type Item[T any] struct {
	Value    T
	priority int
	index    int
}

func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		less: less,
	}
}

func (pq PriorityQueue[T]) Len() int { return len(pq.items) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// First compare by priority (lower value means higher priority)
	if pq.items[i].priority != pq.items[j].priority {
		return pq.items[i].priority < pq.items[j].priority
	}
	// If priorities are equal, use the provided less function
	return pq.less(pq.items[i].Value, pq.items[j].Value)
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(pq.items)
	item := x.(*Item[T])
	item.index = n
	pq.items = append(pq.items, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	item.index = -1 // For safety
	pq.items = old[0 : n-1]
	return item
}

func (pq *PriorityQueue[T]) Enqueue(value T, priority int) {
	heap.Push(pq, &Item[T]{
		Value:    value,
		priority: priority,
	})
}

// Dequeue removes and returns the highest priority item
func (pq *PriorityQueue[T]) Dequeue() (T, bool) {
	if pq.Len() == 0 {
		var zero T
		return zero, false
	}
	item := heap.Pop(pq).(*Item[T])
	return item.Value, true
}
