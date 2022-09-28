package service

import (
	"net"
)

type Service interface {
	Serve() error
	Addr() net.Addr
	Close() error
}
