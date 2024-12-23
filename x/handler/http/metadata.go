package http

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"github.com/jxo-me/netx/x/app"
	"net/http"
	"strings"
	"time"

	"github.com/jxo-me/netx/core/bypass"
	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultRealm      = "gost"
	defaultProxyAgent = "gost/3.0"
)

type metadata struct {
	readTimeout            time.Duration
	keepalive              bool
	compression            bool
	probeResistance        *probeResistance
	enableUDP              bool
	header                 http.Header
	hash                   string
	authBasicRealm         string
	proxyAgent             string
	observePeriod          time.Duration
	limiterRefreshInterval time.Duration

	sniffing                    bool
	sniffingTimeout             time.Duration
	sniffingWebsocket           bool
	sniffingWebsocketSampleRate float64

	certificate *x509.Certificate
	privateKey  crypto.PrivateKey
	alpn        string
	mitmBypass  bypass.IBypass
}

func (h *httpHandler) parseMetadata(md mdata.IMetaData) error {
	h.md.readTimeout = mdutil.GetDuration(md, "readTimeout")
	if h.md.readTimeout <= 0 {
		h.md.readTimeout = 15 * time.Second
	}

	if m := mdutil.GetStringMapString(md, "http.header", "header"); len(m) > 0 {
		hd := http.Header{}
		for k, v := range m {
			hd.Add(k, v)
		}
		h.md.header = hd
	}

	h.md.keepalive = mdutil.GetBool(md, "http.keepalive", "keepalive")
	h.md.compression = mdutil.GetBool(md, "http.compression", "compression")

	if pr := mdutil.GetString(md, "probeResist", "probe_resist"); pr != "" {
		if ss := strings.SplitN(pr, ":", 2); len(ss) == 2 {
			h.md.probeResistance = &probeResistance{
				Type:  ss[0],
				Value: ss[1],
				Knock: mdutil.GetString(md, "knock"),
			}
		}
	}
	h.md.enableUDP = mdutil.GetBool(md, "udp")
	h.md.hash = mdutil.GetString(md, "hash")
	h.md.authBasicRealm = mdutil.GetString(md, "authBasicRealm")

	h.md.observePeriod = mdutil.GetDuration(md, "observePeriod", "observer.observePeriod")
	if h.md.observePeriod == 0 {
		h.md.observePeriod = 5 * time.Second
	}
	if h.md.observePeriod < time.Second {
		h.md.observePeriod = time.Second
	}

	h.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if h.md.limiterRefreshInterval == 0 {
		h.md.limiterRefreshInterval = 30 * time.Second
	}
	if h.md.limiterRefreshInterval < time.Second {
		h.md.limiterRefreshInterval = time.Second
	}

	h.md.proxyAgent = mdutil.GetString(md, "http.proxyAgent", "proxyAgent")
	if h.md.proxyAgent == "" {
		h.md.proxyAgent = defaultProxyAgent
	}

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

	return nil
}

type probeResistance struct {
	Type  string
	Value string
	Knock string
}
