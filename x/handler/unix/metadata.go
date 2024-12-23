package unix

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"github.com/jxo-me/netx/x/app"
	"time"

	"github.com/jxo-me/netx/core/bypass"
	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	readTimeout time.Duration

	sniffing        bool
	sniffingTimeout time.Duration

	certificate *x509.Certificate
	privateKey  crypto.PrivateKey
	alpn        string
	mitmBypass  bypass.IBypass
}

func (h *unixHandler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.readTimeout = mdutil.GetDuration(md, "readTimeout")
	if h.md.readTimeout <= 0 {
		h.md.readTimeout = 15 * time.Second
	}

	h.md.sniffing = mdutil.GetBool(md, "sniffing")
	h.md.sniffingTimeout = mdutil.GetDuration(md, "sniffing.timeout")

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
	return
}
