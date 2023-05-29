package parsing

import "context"

type rpcCredentials struct {
	token string
}

func (c *rpcCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": c.token,
	}, nil
}

func (c *rpcCredentials) RequireTransportSecurity() bool {
	return false
}
