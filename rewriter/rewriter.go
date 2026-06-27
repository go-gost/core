package rewriter

import "context"

type RewriteOptions struct {
	Metadata any
}

type RewriteOption func(opts *RewriteOptions)

func MetadataRewriteOption(md any) RewriteOption {
	return func(opts *RewriteOptions) {
		opts.Metadata = md
	}
}

type Rewriter interface {
	Rewrite(ctx context.Context, b []byte, opts ...RewriteOption) ([]byte, error)
}
