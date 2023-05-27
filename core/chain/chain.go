package chain

import (
	"context"
)

type IChainer interface {
	Route(ctx context.Context, network, address string) IRoute
}
