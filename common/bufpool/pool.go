// Package bufpool provides a tiered sync.Pool for byte buffers, ranging from
// 128 bytes to 65 KB. Using pooled buffers reduces GC pressure in high-throughput
// proxy scenarios where buffers are allocated and discarded frequently.
package bufpool

import (
	"sync"
)

var (
	pools = []struct {
		size int
		pool sync.Pool
	}{
		{
			size: 128,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 128)
					return b
				},
			},
		},
		{
			size: 512,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 512)
					return b
				},
			},
		},
		{
			size: 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 1024)
					return b
				},
			},
		},
		{
			size: 2048,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 2048)
					return b
				},
			},
		},
		{
			size: 4096,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 4096)
					return b
				},
			},
		},
		{
			size: 8192,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 8192)
					return b
				},
			},
		},
		{
			size: 16 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 16*1024)
					return b
				},
			},
		},
		{
			size: 32 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 32*1024)
					return b
				},
			},
		},
		{
			size: 64 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 64*1024)
					return b
				},
			},
		},
		{
			size: 65 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 65*1024)
					return b
				},
			},
		},
	}
)

// Get returns a byte buffer of at least the requested size.
// The buffer is sourced from the smallest pool tier that satisfies the size,
// or allocated fresh if no pool tier is large enough.
func Get(size int) []byte {
	for i := range pools {
		if size <= pools[i].size {
			b := pools[i].pool.Get().([]byte)
			return b[:size]
		}
	}
	b := make([]byte, size)
	return b
}

// Put returns a byte buffer to the appropriate pool tier based on its capacity.
func Put(b []byte) {
	for i := range pools {
		if cap(b) == pools[i].size {
			pools[i].pool.Put(b)
		}
	}
}
