package recorder

import (
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/recorder"
)

type tcpRecorderOptions struct {
	timeout time.Duration
	log     logger.ILogger
}

type TCPRecorderOption func(opts *tcpRecorderOptions)

func TimeoutTCPRecorderOption(timeout time.Duration) TCPRecorderOption {
	return func(opts *tcpRecorderOptions) {
		opts.timeout = timeout
	}
}

func LogTCPRecorderOption(log logger.ILogger) TCPRecorderOption {
	return func(opts *tcpRecorderOptions) {
		opts.log = log
	}
}

type tcpRecorder struct {
	addr   string
	dialer *net.Dialer
	log    logger.ILogger
}

// TCPRecorder records data to TCP service.
func TCPRecorder(addr string, opts ...TCPRecorderOption) recorder.IRecorder {
	var options tcpRecorderOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &tcpRecorder{
		addr: addr,
		dialer: &net.Dialer{
			Timeout: options.timeout,
		},
		log: options.log,
	}
}

func (r *tcpRecorder) Record(ctx context.Context, b []byte, opts ...recorder.RecordOption) error {
	c, err := r.dialer.DialContext(ctx, "tcp", r.addr)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Write(b)
	return err
}
