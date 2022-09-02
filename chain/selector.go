package chain

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-gost/core/metadata"
)

// default options for FailFilter
const (
	DefaultFailTimeout = 30 * time.Second
)

var (
	DefaultNodeSelector = NewSelector(
		RoundRobinStrategy[*Node](),
		// FailFilter[*Node](1, DefaultFailTimeout),
	)
	DefaultChainSelector = NewSelector(
		RoundRobinStrategy[SelectableChainer](),
		// FailFilter[SelectableChainer](1, DefaultFailTimeout),
	)
)

type Selectable interface {
	Marker() Marker
	Metadata() metadata.Metadata
}

type Selector[T any] interface {
	Select(...T) T
}

type selector[T Selectable] struct {
	strategy Strategy[T]
	filters  []Filter[T]
}

func NewSelector[T Selectable](strategy Strategy[T], filters ...Filter[T]) Selector[T] {
	return &selector[T]{
		filters:  filters,
		strategy: strategy,
	}
}

func (s *selector[T]) Select(vs ...T) (v T) {
	for _, filter := range s.filters {
		vs = filter.Filter(vs...)
	}
	if len(vs) == 0 {
		return
	}
	return s.strategy.Apply(vs...)
}

type Strategy[T Selectable] interface {
	Apply(...T) T
}

type roundRobinStrategy[T Selectable] struct {
	counter uint64
}

// RoundRobinStrategy is a strategy for node selector.
// The node will be selected by round-robin algorithm.
func RoundRobinStrategy[T Selectable]() Strategy[T] {
	return &roundRobinStrategy[T]{}
}

func (s *roundRobinStrategy[T]) Apply(vs ...T) (v T) {
	if len(vs) == 0 {
		return
	}

	n := atomic.AddUint64(&s.counter, 1) - 1
	return vs[int(n%uint64(len(vs)))]
}

type randomStrategy[T Selectable] struct {
	rand *rand.Rand
	mux  sync.Mutex
}

// RandomStrategy is a strategy for node selector.
// The node will be selected randomly.
func RandomStrategy[T Selectable]() Strategy[T] {
	return &randomStrategy[T]{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *randomStrategy[T]) Apply(vs ...T) (v T) {
	if len(vs) == 0 {
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	r := s.rand.Int()

	return vs[r%len(vs)]
}

type fifoStrategy[T Selectable] struct{}

// FIFOStrategy is a strategy for node selector.
// The node will be selected from first to last,
// and will stick to the selected node until it is failed.
func FIFOStrategy[T Selectable]() Strategy[T] {
	return &fifoStrategy[T]{}
}

// Apply applies the fifo strategy for the nodes.
func (s *fifoStrategy[T]) Apply(vs ...T) (v T) {
	if len(vs) == 0 {
		return
	}
	return vs[0]
}

type Filter[T Selectable] interface {
	Filter(...T) []T
}

type failFilter[T Selectable] struct {
	maxFails    int
	failTimeout time.Duration
}

// FailFilter filters the dead node.
// A node is marked as dead if its failed count is greater than MaxFails.
func FailFilter[T Selectable](maxFails int, timeout time.Duration) Filter[T] {
	return &failFilter[T]{
		maxFails:    maxFails,
		failTimeout: timeout,
	}
}

// Filter filters dead nodes.
func (f *failFilter[T]) Filter(vs ...T) []T {
	maxFails := f.maxFails
	failTimeout := f.failTimeout
	if failTimeout == 0 {
		failTimeout = DefaultFailTimeout
	}

	if len(vs) <= 1 || maxFails <= 0 {
		return vs
	}
	var l []T
	for _, v := range vs {
		if marker := v.Marker(); marker != nil {
			if marker.Count() < int64(maxFails) ||
				time.Since(marker.Time()) >= failTimeout {
				l = append(l, v)
			}
		} else {
			l = append(l, v)
		}
	}
	return l
}

type Marker interface {
	Time() time.Time
	Count() int64
	Mark()
	Reset()
}

type failMarker struct {
	failTime  int64
	failCount int64
}

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
