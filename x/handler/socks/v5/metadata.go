package v5

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"github.com/jxo-me/netx/x/app"
	"time"

	"github.com/jxo-me/netx/core/bypass"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/x/internal/util/mux"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	readTimeout       time.Duration
	noTLS             bool
	enableBind        bool
	enableUDP         bool
	compatibilityMode bool
	hash              string
	muxCfg            *mux.Config

	observerPeriod       time.Duration
	observerResetTraffic bool

	sniffing                    bool
	sniffingTimeout             time.Duration
	sniffingWebsocket           bool
	sniffingWebsocketSampleRate float64

	certificate *x509.Certificate
	privateKey  crypto.PrivateKey
	alpn        string
	mitmBypass  bypass.IBypass

	limiterRefreshInterval time.Duration
	limiterCleanupInterval time.Duration
}

func (h *socks5Handler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.readTimeout = mdutil.GetDuration(md, "readTimeout")
	if h.md.readTimeout <= 0 {
		h.md.readTimeout = 15 * time.Second
	}

	h.md.noTLS = mdutil.GetBool(md, "notls")
	h.md.enableBind = mdutil.GetBool(md, "bind")
	h.md.enableUDP = mdutil.GetBool(md, "udp")

	h.md.compatibilityMode = mdutil.GetBool(md, "comp")
	h.md.hash = mdutil.GetString(md, "hash")

	h.md.muxCfg = &mux.Config{
		Version:           mdutil.GetInt(md, "mux.version"),
		KeepAliveInterval: mdutil.GetDuration(md, "mux.keepaliveInterval"),
		KeepAliveDisabled: mdutil.GetBool(md, "mux.keepaliveDisabled"),
		KeepAliveTimeout:  mdutil.GetDuration(md, "mux.keepaliveTimeout"),
		MaxFrameSize:      mdutil.GetInt(md, "mux.maxFrameSize"),
		MaxReceiveBuffer:  mdutil.GetInt(md, "mux.maxReceiveBuffer"),
		MaxStreamBuffer:   mdutil.GetInt(md, "mux.maxStreamBuffer"),
	}

	h.md.observerPeriod = mdutil.GetDuration(md, "observePeriod", "observer.period", "observer.observePeriod")
	if h.md.observerPeriod == 0 {
		h.md.observerPeriod = 5 * time.Second
	}
	if h.md.observerPeriod < time.Second {
		h.md.observerPeriod = time.Second
	}
	h.md.observerResetTraffic = mdutil.GetBool(md, "observer.resetTraffic")

	h.md.sniffing = mdutil.GetBool(md, "sniffing")
	h.md.sniffingTimeout = mdutil.GetDuration(md, "sniffing.timeout")
	h.md.sniffingWebsocket = mdutil.GetBool(md, "sniffing.websocket")
	h.md.sniffingWebsocketSampleRate = mdutil.GetFloat(md, "sniffing.websocket.sampleRate")

	certFile := mdutil.GetString(md, "mitm.certFile", "mitm.caCertFile")
	keyFile := mdutil.GetString(md, "mitm.keyFile", "mitm.caKeyFile")
	if certFile != "" && keyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		h.md.certificate, err = x509.ParseCertificate(tlsCert.Certificate[0])
		if err != nil {
			return err
		}
		h.md.privateKey = tlsCert.PrivateKey
	}
	h.md.alpn = mdutil.GetString(md, "mitm.alpn")
	h.md.mitmBypass = app.Runtime.BypassRegistry().Get(mdutil.GetString(md, "mitm.bypass"))

	h.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	h.md.limiterCleanupInterval = mdutil.GetDuration(md, "limiter.cleanupInterval")

	return nil
}
