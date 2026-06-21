// Package recorder defines the Recorder interface for capturing traffic
// data — hex dumps, HTTP body recording, and directional traffic logging.
package recorder

import (
	"context"
)

// RecordOptions holds runtime parameters for a Record call.
type RecordOptions struct {
	// Metadata is contextual data attached to the record operation.
	Metadata any
}

// RecordOption is a functional option for configuring RecordOptions.
type RecordOption func(opts *RecordOptions)

// MetadataRecordOption sets the metadata for the record operation.
func MetadataRecordOption(md any) RecordOption {
	return func(opts *RecordOptions) {
		opts.Metadata = md
	}
}

// Recorder captures traffic data. Implementations may write hex dumps,
// log HTTP request/response bodies, or forward data to external collectors.
type Recorder interface {
	// Record processes a buffer of traffic data.
	Record(ctx context.Context, b []byte, opts ...RecordOption) error
}

// RecorderObject bundles a Recorder with its name and configuration.
// It is used in RouterOptions and Handler options to attach recording
// to specific points in the request path.
type RecorderObject struct {
	// Recorder is the recorder implementation.
	Recorder Recorder
	// Record is the recorder's registered name.
	Record string
	// Options configures what data is recorded.
	Options *Options
	// Metadata is arbitrary metadata attached to the recorder,
	// forwarded to plugin recorders on every Record call.
	Metadata map[string]any
}

// Options configures a Recorder's behavior.
type Options struct {
	// Direction indicates the traffic direction (true = client-to-server).
	Direction bool
	// TimestampFormat is the Go time format string for timestamps.
	TimestampFormat string
	// Hexdump enables hexadecimal dump output.
	Hexdump bool
	// HTTPBody enables HTTP request/response body capture.
	HTTPBody bool
	// MaxBodySize is the maximum body size to record.
	MaxBodySize int
}

// Recorder event key constants used as metadata keys in record operations.
const (
	RecorderServiceClientAddress          = "recorder.service.client.address"
	RecorderServiceRouterDialAddress      = "recorder.service.router.dial.address"
	RecorderServiceRouterDialAddressError = "recorder.service.router.dial.address.error"
)
