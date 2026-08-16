package region

import "context"

type Enricher interface {
	Enrich(ctx context.Context, address *Address) error
}
