package tls

import (
	"time"

	mdata "github.com/go-gost/core/metadata"
)

type metadata struct {
	handshakeTimeout time.Duration
}

func (d *tlsDialer) parseMetadata(md mdata.Metadata) (err error) {
	const (
		handshakeTimeout = "handshakeTimeout"
	)

	d.md.handshakeTimeout = mdata.GetDuration(md, handshakeTimeout)

	return
}
