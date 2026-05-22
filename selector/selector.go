// Package selector defines generic interfaces for node selection with
// support for strategies, filters, and failure marking.
package selector

import (
	"context"
	"sync/atomic"
	"time"
)

// Selector picks one item from a list using a configured strategy.
type Selector[T any] interface {
	Select(context.Context, ...T) T
}

// Strategy applies a selection algorithm (e.g. round-robin, random,
// weighted) to choose one item from the list.
type Strategy[T any] interface {
	Apply(context.Context, ...T) T
}

// Filter narrows a list of items to those that pass the filter criteria.
// Filters are applied before the selection strategy.
type Filter[T any] interface {
	Filter(context.Context, ...T) []T
}

// Markable exposes a Marker for health tracking. Nodes implement this
// so that selectors can temporarily skip failed nodes.
type Markable interface {
	Marker() Marker
}

// Marker tracks failure state for a selectable item. When a node fails,
// it is marked; the selection strategy can then deprioritize or skip
// marked nodes until they recover.
type Marker interface {
	// Time returns the timestamp of the last failure.
	Time() time.Time
	// Count returns the number of consecutive failures.
	Count() int64
	// Mark records a failure.
	Mark()
	// Reset clears the failure state.
	Reset()
}

// failMarker is the default Marker implementation using atomic operations
// for thread safety. It tracks the time and count of the most recent failure.
type failMarker struct {
	failTime  int64
	failCount int64
}

// NewFailMarker creates a new Marker for tracking node failures.
func NewFailMarker() Marker {
	return &failMarker{}
}

func (m *failMarker) Time() time.Time {
	if m == nil {
		return time.Time{}
	}

	return time.Unix(atomic.LoadInt64(&m.failTime), 0)
}

func (m *failMarker) Count() int64 {
	if m == nil {
		return 0
	}

	return atomic.LoadInt64(&m.failCount)
}

func (m *failMarker) Mark() {
	if m == nil {
		return
	}

	atomic.AddInt64(&m.failCount, 1)
	atomic.StoreInt64(&m.failTime, time.Now().Unix())
}

func (m *failMarker) Reset() {
	if m == nil {
		return
	}

	atomic.StoreInt64(&m.failCount, 0)
}
