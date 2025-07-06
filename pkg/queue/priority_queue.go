package queue

import (
	"container/heap"
)

// PriorityQueue implements a generic priority queue
type PriorityQueue[T any] struct {
	items []*Item[T]
	less  func(a, b T) bool // Comparison function
}

// Item represents an element in the priority queue
type Item[T any] struct {
	Value    T   // The generic value
	priority int // The priority of the item
	index    int // The index in the heap
}

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		less: less,
	}
}

// Implement heap.Interface
func (pq PriorityQueue[T]) Len() int { return len(pq.items) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// First compare by priority (higher priority first)
	if pq.items[i].priority != pq.items[j].priority {
		return pq.items[i].priority > pq.items[j].priority
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

// Enqueue adds an item to the queue
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

// Example usage:
// func main() {
// 	// Create queue with string values
// 	pq := NewPriorityQueue[string](func(a, b string) bool {
// 		return a < b // Alphabetical order for equal priorities
// 	})

// 	// Add items
// 	pq.Enqueue("apple", 3)
// 	pq.Enqueue("banana", 2)
// 	pq.Enqueue("orange", 3) // Same priority as apple
// 	pq.Enqueue("pear", 1)

// 	// Process items
// 	for {
// 		item, ok := pq.Dequeue()
// 		if !ok {
// 			break
// 		}
// 		fmt.Println(item)
// 	}
// 	// Output:
// 	// apple
// 	// orange
// 	// banana
// 	// pear
// }
