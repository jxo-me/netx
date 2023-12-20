package serial

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hop"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/recorder"
	xnet "github.com/jxo-me/netx/x/internal/net"
	serial "github.com/jxo-me/netx/x/internal/util/serial"
	xrecorder "github.com/jxo-me/netx/x/recorder"
)

type serialHandler struct {
	hop      hop.IHop
	router   *chain.Router
	md       metadata
	options  handler.Options
	recorder recorder.RecorderObject
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &serialHandler{
		options: options,
	}
}

func (h *serialHandler) Init(md md.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	h.router = h.options.Router
	if h.router == nil {
		h.router = chain.NewRouter(chain.LoggerRouterOption(h.options.Logger))
	}
	if opts := h.router.Options(); opts != nil {
		for _, ro := range opts.Recorders {
			if ro.Record == xrecorder.RecorderServiceHandlerSerial {
				h.recorder = ro
				break
			}
		}
	}

	return
}

// Forward implements handler.Forwarder.
func (h *serialHandler) Forward(hop hop.IHop) {
	h.hop = hop
}

func (h *serialHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	defer conn.Close()

	log := h.options.Logger

	log = log.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
	})

	conn = &recorderConn{
		Conn:     conn,
		recorder: h.recorder,
	}

	if h.hop != nil {
		target := h.hop.Select(ctx)
		if target == nil {
			err := errors.New("target not available")
			log.Error(err)
			return err
		}
		log = log.WithFields(map[string]any{
			"node": target.Name,
			"dst":  target.Addr,
		})
		return h.forwardSerial(ctx, conn, target, log)
	}

	cc, err := h.router.Dial(ctx, "tcp", "@")
	if err != nil {
		log.Error(err)
		return err
	}
	defer cc.Close()

	t := time.Now()
	log.Infof("%s <-> %s", conn.LocalAddr(), "@")
	xnet.Transport(conn, cc)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", conn.LocalAddr(), "@")

	return nil
}

func (h *serialHandler) forwardSerial(ctx context.Context, conn net.Conn, target *chain.Node, log logger.ILogger) (err error) {
	log.Debugf("%s >> %s", conn.LocalAddr(), target.Addr)
	var port io.ReadWriteCloser

	cfg := serial.ParseConfigFromAddr(conn.LocalAddr().String())
	cfg.Name = target.Addr

	if opts := h.router.Options(); opts != nil && opts.Chain != nil {
		port, err = h.router.Dial(ctx, "serial", serial.AddrFromConfig(cfg))
	} else {
		cfg.ReadTimeout = h.md.timeout
		port, err = serial.OpenPort(cfg)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	defer port.Close()

	t := time.Now()
	log.Infof("%s <-> %s", conn.LocalAddr(), target.Addr)
	xnet.Transport(conn, port)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", conn.LocalAddr(), target.Addr)

	return nil
}
