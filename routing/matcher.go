package routing

import (
	"net"
	"net/http"
	"net/url"
)

type Request struct {
	ClientIP net.IP
	Host     string
	Protocol string
	Method   string
	Path     string
	Query    url.Values
	Header   http.Header
}

type Matcher interface {
	Match(*Request) bool
}
