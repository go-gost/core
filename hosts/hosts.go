package hosts

import (
	"net"
)

// HostMapper is a mapping from hostname to IP.
type HostMapper interface {
	Lookup(network, host string) ([]net.IP, bool)
}
