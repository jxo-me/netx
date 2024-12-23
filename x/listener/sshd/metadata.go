package ssh

import (
	"fmt"
	"os"
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	ssh_util "github.com/jxo-me/netx/x/internal/util/ssh"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
	"github.com/mitchellh/go-homedir"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh"
)

const (
	defaultBacklog = 128
)

type metadata struct {
	signer                 ssh.Signer
	authorizedKeys         map[string]bool
	backlog                int
	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *sshdListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		authorizedKeys = "authorizedKeys"
		privateKeyFile = "privateKeyFile"
		passphrase     = "passphrase"
		backlog        = "backlog"
	)

	if key := mdutil.GetString(md, privateKeyFile); key != "" {
		key, err = homedir.Expand(key)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(key)
		if err != nil {
			return err
		}

		var pp string
		if mdutil.GetBool(md, "passphraseFromKeyring") {
			pp, err = keyring.Get(fmt.Sprintf("SSH %s", key), l.options.Auth.Username())
			if err != nil {
				return fmt.Errorf("unable to get secret(%s) from keyring: %w", key, err)
			}
		} else {
			pp = mdutil.GetString(md, passphrase)
		}
		if pp == "" {
			l.md.signer, err = ssh.ParsePrivateKey(data)
		} else {
			l.md.signer, err = ssh.ParsePrivateKeyWithPassphrase(data, []byte(pp))
		}
		if err != nil {
			return err
		}
	}
	if l.md.signer == nil {
		signer, err := ssh.NewSignerFromKey(l.options.TLSConfig.Certificates[0].PrivateKey)
		if err != nil {
			return err
		}
		l.md.signer = signer
	}

	if name := mdutil.GetString(md, authorizedKeys); name != "" {
		m, err := ssh_util.ParseAuthorizedKeysFile(name)
		if err != nil {
			return err
		}
		l.md.authorizedKeys = m
	}

	l.md.backlog = mdutil.GetInt(md, backlog)
	if l.md.backlog <= 0 {
		l.md.backlog = defaultBacklog
	}

	l.md.mptcp = mdutil.GetBool(md, "mptcp")

	l.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if l.md.limiterRefreshInterval == 0 {
		l.md.limiterRefreshInterval = 30 * time.Second
	}
	if l.md.limiterRefreshInterval < time.Second {
		l.md.limiterRefreshInterval = time.Second
	}

	return
}
