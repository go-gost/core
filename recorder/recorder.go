package recorder

import (
	"context"
)

type RecordOptions struct{}
type RecordOption func(opts *RecordOptions)

type Recorder interface {
	Record(ctx context.Context, b []byte, opts ...RecordOption) error
}

type RecorderObject struct {
	Recorder Recorder
	Record   string
	Options  *Options
}

type Options struct {
	Direction       bool
	TimestampFormat string
	Hexdump         bool
}

const (
	RecorderServiceClientAddress          = "recorder.service.client.address"
	RecorderServiceRouterDialAddress      = "recorder.service.router.dial.address"
	RecorderServiceRouterDialAddressError = "recorder.service.router.dial.address.error"
)
