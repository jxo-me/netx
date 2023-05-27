package ingress

import "context"

type IIngress interface {
	Get(ctx context.Context, host string) string
}
