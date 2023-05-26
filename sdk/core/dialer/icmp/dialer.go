package quic

import (
	"context"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"golang.org/x/net/icmp"
	"github.com/jxo-me/netx/sdk/core/dialer"
	icmp_pkg "github.com/jxo-me/netx/sdk/internal/util/icmp"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
)

type icmpDialer struct {
	sessions     map[string]*quicSession
	sessionMutex sync.Mutex
	logger       logger.ILogger
	md           metadata
	options      dialer.Options
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := dialer.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &icmpDialer{
		sessions: make(map[string]*quicSession),
		logger:   options.Logger,
		options:  options,
	}
}

func (d *icmpDialer) Init(md md.IMetaData) (err error) {
	if err = d.parseMetadata(md); err != nil {
		return
	}

	return nil
}

func (d *icmpDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (conn net.Conn, err error) {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "0")
	}

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	d.sessionMutex.Lock()
	defer d.sessionMutex.Unlock()

	session, ok := d.sessions[addr]
	if !ok {
		options := &dialer.DialOptions{}
		for _, opt := range opts {
			opt(options)
		}

		var pc net.PacketConn
		pc, err = icmp.ListenPacket("ip4:icmp", "")
		if err != nil {
			return
		}

		id := raddr.Port
		if id == 0 {
			id = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(math.MaxUint16) + 1
			raddr.Port = id
		}
		pc = icmp_pkg.ClientConn(pc, id)

		session, err = d.initSession(ctx, raddr, pc)
		if err != nil {
			d.logger.Error(err)
			pc.Close()
			return nil, err
		}

		d.sessions[addr] = session
	}

	conn, err = session.GetConn()
	if err != nil {
		session.Close()
		delete(d.sessions, addr)
		return nil, err
	}

	return
}

func (d *icmpDialer) initSession(ctx context.Context, addr net.Addr, conn net.PacketConn) (*quicSession, error) {
	quicConfig := &quic.Config{
		KeepAlivePeriod:      d.md.keepAlivePeriod,
		HandshakeIdleTimeout: d.md.handshakeTimeout,
		MaxIdleTimeout:       d.md.maxIdleTimeout,
		Versions: []quic.VersionNumber{
			quic.Version1,
			quic.VersionDraft29,
		},
	}

	tlsCfg := d.options.TLSConfig
	tlsCfg.NextProtos = []string{"http/3", "quic/v1"}

	session, err := quic.DialEarlyContext(ctx, conn, addr, addr.String(), tlsCfg, quicConfig)
	if err != nil {
		return nil, err
	}
	return &quicSession{session: session}, nil
}

// Multiplex implements dialer.Multiplexer interface.
func (d *icmpDialer) Multiplex() bool {
	return true
}
