package metrics

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/jxo-me/netx/core/handler"
	md "github.com/jxo-me/netx/core/metadata"
	xmetrics "github.com/jxo-me/netx/x/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsHandler struct {
	handler http.Handler
	mux     *http.ServeMux
	md      metadata
	options handler.Options
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &metricsHandler{
		options: options,
	}
}

func (h *metricsHandler) Init(md md.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	xmetrics.Enable(true)

	h.handler = promhttp.Handler()

	mux := http.NewServeMux()
	mux.Handle(h.md.path, http.HandlerFunc(h.handleFunc))
	h.mux = mux

	return
}

func (h *metricsHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	l := &singleConnListener{
		conn: make(chan net.Conn, 1),
	}
	l.send(conn)

	s := http.Server{
		Handler: h.mux,
	}
	s.Serve(l)

	return s.Shutdown(ctx)
}

func (h *metricsHandler) Close() error {
	return nil
}

func (h *metricsHandler) handleFunc(w http.ResponseWriter, r *http.Request) {
	if auther := h.options.Auther; auther != nil {
		u, p, _ := r.BasicAuth()
		if _, ok := auther.Authenticate(r.Context(), u, p); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	log := h.options.Logger
	start := time.Now()

	h.handler.ServeHTTP(w, r)

	log = log.WithFields(map[string]any{
		"remote":   r.RemoteAddr,
		"duration": time.Since(start),
	})
	log.Infof("%s %s", r.Method, r.RequestURI)
}

type singleConnListener struct {
	conn chan net.Conn
}

func (l *singleConnListener) Accept() (net.Conn, error) {
	if conn, ok := <-l.conn; ok {
		return conn, nil
	}
	return nil, net.ErrClosed
}

func (l *singleConnListener) Close() error {
	return nil
}

func (l *singleConnListener) Addr() net.Addr {
	return &net.TCPAddr{}
}

func (l *singleConnListener) send(conn net.Conn) {
	l.conn <- conn
	close(l.conn)
}
