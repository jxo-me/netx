package metrics

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/jxo-me/netx/core/handler"
	md "github.com/jxo-me/netx/core/metadata"
	xmetrics "github.com/jxo-me/netx/x/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsHandler struct {
	handler http.Handler
	server  *http.Server
	ln      *singleConnListener
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

	xmetrics.Init(xmetrics.NewMetrics())
	h.handler = promhttp.Handler()

	mux := http.NewServeMux()
	mux.Handle(h.md.path, http.HandlerFunc(h.handleFunc))
	h.server = &http.Server{
		Handler: mux,
	}

	h.ln = &singleConnListener{
		conn: make(chan net.Conn),
		done: make(chan struct{}),
	}
	go h.server.Serve(h.ln)

	return
}

func (h *metricsHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	h.ln.send(conn)

	return nil
}

func (h *metricsHandler) Close() error {
	return h.server.Close()
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
	log.Debugf("%s %s", r.Method, r.RequestURI)
}

type singleConnListener struct {
	conn chan net.Conn
	addr net.Addr
	done chan struct{}
	mu   sync.Mutex
}

func (l *singleConnListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.conn:
		return conn, nil

	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *singleConnListener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	select {
	case <-l.done:
	default:
		close(l.done)
	}

	return nil
}

func (l *singleConnListener) Addr() net.Addr {
	return l.addr
}

func (l *singleConnListener) send(conn net.Conn) {
	select {
	case l.conn <- conn:
	case <-l.done:
		return
	}
}
