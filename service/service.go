package service

import (
	"context"
	"net"
	"time"

	"github.com/go-gost/core/admission"
	"github.com/go-gost/core/handler"
	"github.com/go-gost/core/listener"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/metrics"
)

type options struct {
	admission admission.Admission
	logger    logger.Logger
}

type Option func(opts *options)

func AdmissionOption(admission admission.Admission) Option {
	return func(opts *options) {
		opts.admission = admission
	}
}

func LoggerOption(logger logger.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type Service interface {
	Serve() error
	Addr() net.Addr
	Close() error
}

type service struct {
	name     string
	listener listener.Listener
	handler  handler.Handler
	options  options
}

func NewService(name string, ln listener.Listener, h handler.Handler, opts ...Option) Service {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	return &service{
		name:     name,
		listener: ln,
		handler:  h,
		options:  options,
	}
}

func (s *service) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *service) Close() error {
	return s.listener.Close()
}

func (s *service) Serve() error {
	if v := metrics.GetGauge(
		metrics.MetricServicesGauge,
		metrics.Labels{}); v != nil {
		v.Inc()
		defer v.Dec()
	}

	var tempDelay time.Duration
	for {
		conn, e := s.listener.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 1 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 5 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.options.logger.Warnf("accept: %v, retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			s.options.logger.Errorf("accept: %v", e)
			return e
		}
		tempDelay = 0

		if s.options.admission != nil &&
			!s.options.admission.Admit(conn.RemoteAddr().String()) {
			conn.Close()
			continue
		}

		go func() {
			if v := metrics.GetCounter(metrics.MetricServiceRequestsCounter,
				metrics.Labels{"service": s.name}); v != nil {
				v.Inc()
			}

			if v := metrics.GetGauge(metrics.MetricServiceRequestsInFlightGauge,
				metrics.Labels{"service": s.name}); v != nil {
				v.Inc()
				defer v.Dec()
			}

			start := time.Now()
			if v := metrics.GetObserver(metrics.MetricServiceRequestsDurationObserver,
				metrics.Labels{"service": s.name}); v != nil {
				defer func() {
					v.Observe(float64(time.Since(start).Seconds()))
				}()
			}

			if err := s.handler.Handle(
				context.Background(),
				conn,
			); err != nil {
				s.options.logger.Error(err)
				if v := metrics.GetCounter(metrics.MetricServiceHandlerErrorsCounter,
					metrics.Labels{"service": s.name}); v != nil {
					v.Inc()
				}
			}
		}()
	}
}
