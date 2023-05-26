package auto

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/sdk/core/handler"
	netpkg "github.com/jxo-me/netx/sdk/internal/net"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	"github.com/jxo-me/netx/sdk/gosocks4"
	"github.com/jxo-me/netx/sdk/gosocks5"
)

type AutoHandler struct {
	httpHandler   handler.IHandler
	socks4Handler handler.IHandler
	socks5Handler handler.IHandler
	options       handler.Options
}

func NewHandler(opts ...handler.Option) *AutoHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	h := &AutoHandler{
		options: options,
	}

	return h
}

func (h *AutoHandler) SetHttpHandler(handle handler.IHandler) {
	h.httpHandler = handle
}

func (h *AutoHandler) SetSocks4Handler(handle handler.IHandler) {
	h.socks4Handler = handle
}

func (h *AutoHandler) SetSocks5Handler(handle handler.IHandler) {
	h.socks5Handler = handle
}

func (h *AutoHandler) Init(md md.IMetaData) error {
	if h.httpHandler != nil {
		if err := h.httpHandler.Init(md); err != nil {
			return err
		}
	}
	if h.socks4Handler != nil {
		if err := h.socks4Handler.Init(md); err != nil {
			return err
		}
	}
	if h.socks5Handler != nil {
		if err := h.socks5Handler.Init(md); err != nil {
			return err
		}
	}

	return nil
}

func (h *AutoHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	log := h.options.Logger.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
	})

	if log.IsLevelEnabled(logger.DebugLevel) {
		start := time.Now()
		log.Debugf("%s <> %s", conn.RemoteAddr(), conn.LocalAddr())
		defer func() {
			log.WithFields(map[string]any{
				"duration": time.Since(start),
			}).Debugf("%s >< %s", conn.RemoteAddr(), conn.LocalAddr())
		}()
	}

	br := bufio.NewReader(conn)
	b, err := br.Peek(1)
	if err != nil {
		log.Error(err)
		conn.Close()
		return err
	}

	conn = netpkg.NewBufferReaderConn(conn, br)
	switch b[0] {
	case gosocks4.Ver4: // socks4
		if h.socks4Handler != nil {
			return h.socks4Handler.Handle(ctx, conn)
		}
	case gosocks5.Ver5: // socks5
		if h.socks5Handler != nil {
			return h.socks5Handler.Handle(ctx, conn)
		}
	default: // http
		if h.httpHandler != nil {
			return h.httpHandler.Handle(ctx, conn)
		}
	}
	return nil
}
