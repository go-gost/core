package udp

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"github.com/go-gost/core/common/bufpool"
)

// conn is a server side connection for UDP client peer, it implements net.Conn and net.PacketConn.
type conn struct {
	net.PacketConn
	localAddr  net.Addr
	remoteAddr net.Addr
	rc         chan []byte // data receive queue
	idle       int32       // indicate the connection is idle
	closed     chan struct{}
	closeMutex sync.Mutex
	keepAlive  bool
}

func newConn(c net.PacketConn, laddr, remoteAddr net.Addr, queueSize int, keepAlive bool) *conn {
	return &conn{
		PacketConn: c,
		localAddr:  laddr,
		remoteAddr: remoteAddr,
		rc:         make(chan []byte, queueSize),
		closed:     make(chan struct{}),
		keepAlive:  keepAlive,
	}
}

func (c *conn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	select {
	case bb := <-c.rc:
		n = copy(b, bb)
		c.SetIdle(false)
		bufpool.Put(&bb)

	case <-c.closed:
		err = net.ErrClosed
		return
	}

	addr = c.remoteAddr

	return
}

func (c *conn) Read(b []byte) (n int, err error) {
	n, _, err = c.ReadFrom(b)
	return
}

func (c *conn) Write(b []byte) (n int, err error) {
	n, err = c.WriteTo(b, c.remoteAddr)
	if !c.keepAlive {
		c.Close()
	}
	return
}

func (c *conn) Close() error {
	c.closeMutex.Lock()
	defer c.closeMutex.Unlock()

	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *conn) IsIdle() bool {
	return atomic.LoadInt32(&c.idle) > 0
}

func (c *conn) SetIdle(idle bool) {
	v := int32(0)
	if idle {
		v = 1
	}
	atomic.StoreInt32(&c.idle, v)
}

func (c *conn) WriteQueue(b []byte) error {
	select {
	case c.rc <- b:
		return nil

	case <-c.closed:
		return net.ErrClosed

	default:
		return errors.New("recv queue is full")
	}
}
