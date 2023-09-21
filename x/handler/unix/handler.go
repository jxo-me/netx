package unix

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	xnet "github.com/jxo-me/netx/x/internal/net"
)

type unixHandler struct {
	hop     chain.IHop
	router  *chain.Router
	md      metadata
	options handler.Options
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &unixHandler{
		options: options,
	}
}

func (h *unixHandler) Init(md md.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	h.router = h.options.Router
	if h.router == nil {
		h.router = chain.NewRouter(chain.LoggerRouterOption(h.options.Logger))
	}

	return
}

// Forward implements handler.Forwarder.
func (h *unixHandler) Forward(hop chain.IHop) {
	h.hop = hop
}

func (h *unixHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	defer conn.Close()

	log := h.options.Logger

	log = log.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
	})

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
		return h.forwardUnix(ctx, conn, target, log)
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

func (h *unixHandler) forwardUnix(ctx context.Context, conn net.Conn, target *chain.Node, log logger.ILogger) (err error) {
	log.Debugf("%s >> %s", conn.LocalAddr(), target.Addr)
	var cc io.ReadWriteCloser

	if opts := h.router.Options(); opts != nil && opts.Chain != nil {
		cc, err = h.router.Dial(ctx, "unix", target.Addr)
	} else {
		cc, err = (&net.Dialer{}).DialContext(ctx, "unix", target.Addr)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	defer cc.Close()

	t := time.Now()
	log.Infof("%s <-> %s", conn.LocalAddr(), target.Addr)
	xnet.Transport(conn, cc)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", conn.LocalAddr(), target.Addr)

	return nil
}
