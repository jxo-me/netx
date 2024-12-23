package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/metrics"
	"github.com/jxo-me/netx/core/observer"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/core/service"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	xnet "github.com/jxo-me/netx/x/internal/net"
	xmetrics "github.com/jxo-me/netx/x/metrics"
	"github.com/rs/xid"
)

type options struct {
	admission     admission.IAdmission
	recorders     []recorder.RecorderObject
	preUp         []string
	postUp        []string
	preDown       []string
	postDown      []string
	stats         *stats.Stats
	observer      observer.IObserver
	observePeriod time.Duration
	logger        logger.ILogger
}

type Option func(opts *options)

func AdmissionOption(admission admission.IAdmission) Option {
	return func(opts *options) {
		opts.admission = admission
	}
}

func RecordersOption(recorders ...recorder.RecorderObject) Option {
	return func(opts *options) {
		opts.recorders = recorders
	}
}

func PreUpOption(cmds []string) Option {
	return func(opts *options) {
		opts.preUp = cmds
	}
}

func PreDownOption(cmds []string) Option {
	return func(opts *options) {
		opts.preDown = cmds
	}
}

func PostUpOption(cmds []string) Option {
	return func(opts *options) {
		opts.postUp = cmds
	}
}

func PostDownOption(cmds []string) Option {
	return func(opts *options) {
		opts.postDown = cmds
	}
}

func StatsOption(stats *stats.Stats) Option {
	return func(opts *options) {
		opts.stats = stats
	}
}

func ObserverOption(observer observer.IObserver) Option {
	return func(opts *options) {
		opts.observer = observer
	}
}

func ObservePeriodOption(period time.Duration) Option {
	return func(opts *options) {
		opts.observePeriod = period
	}
}

func LoggerOption(logger logger.ILogger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type defaultService struct {
	name     string
	listener listener.IListener
	handler  handler.IHandler
	status   *Status
	options  options
}

func NewService(name string, ln listener.IListener, h handler.IHandler, opts ...Option) service.IService {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	s := &defaultService{
		name:     name,
		listener: ln,
		handler:  h,
		options:  options,
		status: &Status{
			createTime: time.Now(),
			events:     make([]Event, 0, MaxEventSize),
			stats:      options.stats,
		},
	}
	s.setState(StateRunning)

	s.execCmds("pre-up", s.options.preUp)

	return s
}

func (s *defaultService) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *defaultService) Serve() error {
	s.execCmds("post-up", s.options.postUp)
	s.setState(StateReady)
	s.status.addEvent(Event{
		Time:    time.Now(),
		Message: fmt.Sprintf("service %s is listening on %s", s.name, s.listener.Addr()),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if s.status.Stats() != nil {
		go s.observeStats(ctx)
	}

	if v := xmetrics.GetGauge(
		xmetrics.MetricServicesGauge,
		metrics.Labels{}); v != nil {
		v.Inc()
		defer v.Dec()
	}

	var tempDelay time.Duration
	for {
		conn, e := s.listener.Accept()
		if e != nil {
			// TODO: remove Temporary checking
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 1 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 5 * time.Second; tempDelay > max {
					tempDelay = max
				}

				s.setState(StateFailed)

				s.options.logger.Warnf("accept: %v, retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			s.setState(StateClosed)
			s.options.logger.Errorf("accept: %v", e)

			return e
		}

		if tempDelay > 0 {
			tempDelay = 0
			s.setState(StateReady)
		}

		clientAddr := conn.RemoteAddr().String()
		if ca, ok := conn.(xnet.ClientAddr); ok {
			if addr := ca.ClientAddr(); addr != nil {
				clientAddr = addr.String()
			}
		}
		clientIP := clientAddr
		if h, _, _ := net.SplitHostPort(clientAddr); h != "" {
			clientIP = h
		}

		sid := xid.New().String()
		ctx := ctxvalue.ContextWithSid(ctx, ctxvalue.Sid(sid))
		ctx = ctxvalue.ContextWithClientAddr(ctx, ctxvalue.ClientAddr(clientAddr))
		ctx = ctxvalue.ContextWithHash(ctx, &ctxvalue.Hash{Source: clientIP})

		log := s.options.logger.WithFields(map[string]any{
			"sid": sid,
		})

		for _, rec := range s.options.recorders {
			if rec.Record == recorder.RecorderServiceClientAddress {
				if err := rec.Recorder.Record(ctx, []byte(clientIP)); err != nil {
					log.Errorf("record %s: %v", rec.Record, err)
				}
				break
			}
		}
		if s.options.admission != nil &&
			!s.options.admission.Admit(ctx, clientAddr) {
			conn.Close()
			log.Debugf("admission: %s is denied", clientAddr)
			continue
		}

		go func() {
			if v := xmetrics.GetCounter(xmetrics.MetricServiceRequestsCounter,
				metrics.Labels{"service": s.name, "client": clientIP}); v != nil {
				v.Inc()
			}

			if v := xmetrics.GetGauge(xmetrics.MetricServiceRequestsInFlightGauge,
				metrics.Labels{"service": s.name, "client": clientIP}); v != nil {
				v.Inc()
				defer v.Dec()
			}

			start := time.Now()
			if v := xmetrics.GetObserver(xmetrics.MetricServiceRequestsDurationObserver,
				metrics.Labels{"service": s.name}); v != nil {
				defer func() {
					v.Observe(float64(time.Since(start).Seconds()))
				}()
			}

			if err := s.handler.Handle(ctx, conn); err != nil {
				log.Error(err)
				if v := xmetrics.GetCounter(xmetrics.MetricServiceHandlerErrorsCounter,
					metrics.Labels{"service": s.name, "client": clientIP}); v != nil {
					v.Inc()
				}
				s.status.stats.Add(stats.KindTotalErrs, 1)
			}
		}()
	}
}

func (s *defaultService) Status() *Status {
	return s.status
}

func (s *defaultService) Close() error {
	s.execCmds("pre-down", s.options.preDown)
	defer s.execCmds("post-down", s.options.postDown)

	if closer, ok := s.handler.(io.Closer); ok {
		closer.Close()
	}
	return s.listener.Close()
}

func (s *defaultService) execCmds(phase string, cmds []string) {
	for _, cmd := range cmds {
		cmd := strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		s.options.logger.Info(cmd)

		if err := exec.Command("/bin/sh", "-c", cmd).Run(); err != nil {
			s.options.logger.Warnf("[%s] %s: %v", phase, cmd, err)
		}
	}
}

func (s *defaultService) setState(state State) {
	s.status.setState(state)

	msg := fmt.Sprintf("service %s is %s", s.name, state)
	s.status.addEvent(Event{
		Time:    time.Now(),
		Message: msg,
	})

	if obs := s.options.observer; obs != nil {
		obs.Observe(context.Background(), []observer.Event{ServiceEvent{
			Kind:    "service",
			Service: s.name,
			State:   state,
			Msg:     msg,
		}})
	}
}

func (s *defaultService) observeStats(ctx context.Context) {
	if s.options.observer == nil {
		return
	}

	d := s.options.observePeriod
	if d < time.Millisecond {
		d = 5 * time.Second
	}

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			st := s.status.Stats()
			if !st.IsUpdated() {
				break
			}
			s.options.observer.Observe(ctx, []observer.Event{
				stats.StatsEvent{
					Kind:         "service",
					Service:      s.name,
					TotalConns:   st.Get(stats.KindTotalConns),
					CurrentConns: st.Get(stats.KindCurrentConns),
					InputBytes:   st.Get(stats.KindInputBytes),
					OutputBytes:  st.Get(stats.KindOutputBytes),
				},
			})
		case <-ctx.Done():
			return
		}
	}
}

type ServiceEvent struct {
	Kind    string
	Service string
	State   State
	Msg     string
}

func (ServiceEvent) Type() observer.EventType {
	return observer.EventStatus
}
