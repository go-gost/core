package chain

import (
	"context"
)

type Chainer interface {
	Route(ctx context.Context, network, address string) Route
}
