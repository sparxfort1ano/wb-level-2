package sort

import (
	"bufio"

	"github.com/sparxfort1ano/wb-level-2/sort/options"
)

// MinHeapItem represents a file as a minimal heap item to be managed in PriorityQueue.
type MinHeapItem struct {
	Line   string
	Reader *bufio.Reader
}

// NewMinHeapItem creates a new instance of MinHeapItem.
func NewMinHeapItem(line string, reader *bufio.Reader) *MinHeapItem {
	return &MinHeapItem{
		Line: line,
		Reader: reader,
	}
}

// PriorityQueue implements [heap.Container] and holds MinHeapItems.
// So that K-Way Merge from sorted temporary files can be implemented.
type PriorityQueue struct {
	Items []*MinHeapItem
	Opts  *options.Options
}

// NewPriorityQueue creates a new instance of PriorityQueue.
func NewPriorityQueue(items []*MinHeapItem, opts *options.Options) *PriorityQueue {
	return &PriorityQueue{
		Items: items,
		Opts: opts,
	}
}

func (pq PriorityQueue) Len() int {
	return len(pq.Items)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq.Opts.Compare(pq.Items[i].Line, pq.Items[j].Line) < 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.Items[i], pq.Items[j] = pq.Items[j], pq.Items[i]
}

func (pq *PriorityQueue) Push(x any) {
	pq.Items = append(pq.Items, x.(*MinHeapItem))
}

func (pq *PriorityQueue) Pop() any {
	n := len(pq.Items)
	x := pq.Items[n-1]

	pq.Items[n-1] = nil

	pq.Items = pq.Items[0 : n-1]
	return x
}
