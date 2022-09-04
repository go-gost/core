package selector

import (
	"context"
	"sync/atomic"
	"time"
)

type Selector[T any] interface {
	Select(context.Context, ...T) T
}

type Strategy[T any] interface {
	Apply(context.Context, ...T) T
}

type Filter[T any] interface {
	Filter(context.Context, ...T) []T
}

type Markable interface {
	Marker() Marker
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
