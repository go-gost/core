package selector

import (
	"sync/atomic"
	"time"

	"github.com/go-gost/core/metadata"
)

type Selectable interface {
	Marker() Marker
	Metadata() metadata.Metadata
}

type Selector[T any] interface {
	Select(...T) T
}

type Strategy[T Selectable] interface {
	Apply(...T) T
}

type Filter[T Selectable] interface {
	Filter(...T) []T
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
