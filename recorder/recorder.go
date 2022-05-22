package recorder

import "context"

type Recorder interface {
	Record(ctx context.Context, b []byte) error
}

type RecorderObject struct {
	Recorder Recorder
	Record   string
}

const (
	RecorderServiceClientAddress          = "recorder.service.client.address"
	RecorderServiceRouterDialAddress      = "recorder.service.router.dial.address"
	RecorderServiceRouterDialAddressError = "recorder.service.router.dial.address.error"
)
