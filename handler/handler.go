package handler

import (
	"context"
	"net"

	"github.com/go-gost/core/chain"
	"github.com/go-gost/core/metadata"
)

type Handler interface {
	Init(metadata.Metadata) error
	Handle(context.Context, net.Conn, ...HandleOption) error
}

type Forwarder interface {
	Forward(chain.Hop)
}
