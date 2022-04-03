package udp

import (
	"net"
	"sync"
	"time"

	"github.com/go-gost/core/common/bufpool"
	"github.com/go-gost/core/logger"
)

type ListenConfig struct {
	Addr           net.Addr
	Backlog        int
	ReadQueueSize  int
	ReadBufferSize int
	TTL            time.Duration
	KeepAlive      bool
	Logger         logger.Logger
}
type listener struct {
	conn     net.PacketConn
	cqueue   chan net.Conn
	connPool *connPool
	mux      sync.Mutex
	closed   chan struct{}
	errChan  chan error
	config   *ListenConfig
}

func NewListener(conn net.PacketConn, cfg *ListenConfig) net.Listener {
	if cfg == nil {
		cfg = &ListenConfig{}
	}

	ln := &listener{
		conn:     conn,
		cqueue:   make(chan net.Conn, cfg.Backlog),
		connPool: newConnPool(cfg.TTL).WithLogger(cfg.Logger),
		closed:   make(chan struct{}),
		errChan:  make(chan error, 1),
		config:   cfg,
	}
	go ln.listenLoop()

	return ln
}

func (ln *listener) Accept() (conn net.Conn, err error) {
	select {
	case conn = <-ln.cqueue:
		return
	case <-ln.closed:
		return nil, net.ErrClosed
	case err = <-ln.errChan:
		if err == nil {
			err = net.ErrClosed
		}
		return
	}
}

func (ln *listener) listenLoop() {
	for {
		select {
		case <-ln.closed:
			return
		default:
		}

		b := bufpool.Get(ln.config.ReadBufferSize)

		n, raddr, err := ln.conn.ReadFrom(*b)
		if err != nil {
			ln.errChan <- err
			close(ln.errChan)
			return
		}

		c := ln.getConn(raddr)
		if c == nil {
			bufpool.Put(b)
			continue
		}

		if err := c.WriteQueue((*b)[:n]); err != nil {
			ln.config.Logger.Warn("data discarded: ", err)
		}
	}
}

func (ln *listener) Addr() net.Addr {
	if ln.config.Addr != nil {
		return ln.config.Addr
	}
	return ln.conn.LocalAddr()
}

func (ln *listener) Close() error {
	select {
	case <-ln.closed:
	default:
		close(ln.closed)
		ln.conn.Close()
		ln.connPool.Close()
	}

	return nil
}

func (ln *listener) getConn(raddr net.Addr) *conn {
	ln.mux.Lock()
	defer ln.mux.Unlock()

	c, ok := ln.connPool.Get(raddr.String())
	if ok {
		return c
	}

	c = newConn(ln.conn, ln.Addr(), raddr, ln.config.ReadQueueSize, ln.config.KeepAlive)
	select {
	case ln.cqueue <- c:
		ln.connPool.Set(raddr.String(), c)
		return c
	default:
		c.Close()
		ln.config.Logger.Warnf("connection queue is full, client %s discarded", raddr)
		return nil
	}
}
