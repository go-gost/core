package dialer

import (
	"context"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/go-gost/core/logger"
)

const (
	DefaultTimeout = 10 * time.Second
)

var (
	DefaultNetDialer = &NetDialer{}
)

type NetDialer struct {
	Interface string
	Mark      int
	Timeout   time.Duration
	DialFunc  func(ctx context.Context, network, addr string) (net.Conn, error)
	Logger    logger.Logger
	deadline  time.Time
}

func (d *NetDialer) dialOnce(ctx context.Context, network, addr, ifceName string, ifAddr net.Addr, log logger.Logger) (net.Conn, error) {
	if ifceName != "" {
		log.Debugf("interface: %s %v/%s", ifceName, ifAddr, network)
	}

	switch network {
	case "udp", "udp4", "udp6":
		if addr == "" {
			var laddr *net.UDPAddr
			if ifAddr != nil {
				laddr, _ = ifAddr.(*net.UDPAddr)
			}

			c, err := net.ListenUDP(network, laddr)
			if err != nil {
				return nil, err
			}
			sc, err := c.SyscallConn()
			if err != nil {
				log.Error(err)
				return nil, err
			}
			err = sc.Control(func(fd uintptr) {
				if ifceName != "" {
					if err := bindDevice(fd, ifceName); err != nil {
						log.Warnf("bind device: %v", err)
					}
				}
				if d.Mark != 0 {
					if err := setMark(fd, d.Mark); err != nil {
						log.Warnf("set mark: %v", err)
					}
				}
			})
			if err != nil {
				log.Error(err)
			}
			return c, nil
		}
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, fmt.Errorf("dial: unsupported network %s", network)
	}
	netd := net.Dialer{
		Deadline:  d.deadline,
		LocalAddr: ifAddr,
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if ifceName != "" {
					if err := bindDevice(fd, ifceName); err != nil {
						log.Warnf("bind device: %v", err)
					}
				}
				if d.Mark != 0 {
					if err := setMark(fd, d.Mark); err != nil {
						log.Warnf("set mark: %v", err)
					}
				}
			})
		},
	}
	return netd.DialContext(ctx, network, addr)
}

func (d *NetDialer) Dial(ctx context.Context, network, addr string) (conn net.Conn, err error) {
	if d == nil {
		d = DefaultNetDialer
	}
	if d.Timeout <= 0 {
		d.Timeout = DefaultTimeout
	}

	if d.DialFunc != nil {
		return d.DialFunc(ctx, network, addr)
	}

	log := d.Logger
	if log == nil {
		log = logger.Default()
	}

	ifces := strings.Split(d.Interface, ",")
	d.deadline = time.Now().Add(d.Timeout)
	for _, ifce := range ifces {
		strict := strings.HasSuffix(ifce, "!")
		ifce = strings.TrimSuffix(ifce, "!")
		var ifceName string
		var ifAddrs []net.Addr
		ifceName, ifAddrs, err = parseInterfaceAddr(ifce, network)
		if err != nil && strict {
			return
		}

		for _, ifAddr := range ifAddrs {
			conn, err = d.dialOnce(ctx, network, addr, ifceName, ifAddr, log)
			if err == nil {
				return
			}

			log.Debugf("dial %s %v@%s failed: %s", network, ifAddr, ifceName, err)

			if strict &&
				!strings.Contains(err.Error(), "no suitable address found") &&
				!strings.Contains(err.Error(), "mismatched local address type") {
				return
			}

			if time.Until(d.deadline) < 0 {
				return
			}
		}
	}

	return
}

func ipToAddr(ip net.IP, network string) (addr net.Addr) {
	port := 0
	switch network {
	case "tcp", "tcp4", "tcp6":
		addr = &net.TCPAddr{IP: ip, Port: port}
		return
	case "udp", "udp4", "udp6":
		addr = &net.UDPAddr{IP: ip, Port: port}
		return
	default:
		addr = &net.IPAddr{IP: ip}
		return
	}
}

func parseInterfaceAddr(ifceName, network string) (ifce string, addr []net.Addr, err error) {
	if ifceName == "" {
		addr = append(addr, nil)
		return
	}

	ip := net.ParseIP(ifceName)
	if ip == nil {
		var ife *net.Interface
		ife, err = net.InterfaceByName(ifceName)
		if err != nil {
			return
		}
		var addrs []net.Addr
		addrs, err = ife.Addrs()
		if err != nil {
			return
		}
		if len(addrs) == 0 {
			err = fmt.Errorf("addr not found for interface %s", ifceName)
			return
		}
		ifce = ifceName
		for _, addr_ := range addrs {
			if ipNet, ok := addr_.(*net.IPNet); ok {
				addr = append(addr, ipToAddr(ipNet.IP, network))
			}
		}
	} else {
		ifce, err = findInterfaceByIP(ip)
		if err != nil {
			return
		}
		addr = []net.Addr{ipToAddr(ip, network)}
	}

	return
}

func findInterfaceByIP(ip net.IP) (string, error) {
	ifces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifce := range ifces {
		addrs, _ := ifce.Addrs()
		if len(addrs) == 0 {
			continue
		}
		for _, addr := range addrs {
			ipAddr, _ := addr.(*net.IPNet)
			if ipAddr == nil {
				continue
			}
			// logger.Default().Infof("%s-%s", ipAddr, ip)
			if ipAddr.IP.Equal(ip) {
				return ifce.Name, nil
			}
		}
	}
	return "", nil
}
