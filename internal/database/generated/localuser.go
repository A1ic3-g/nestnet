package generated

import (
	"golang.org/x/net/context"
)

func (q *Queries) GetName(ctx context.Context) (string, error) {
	return q.getName(ctx)
}

func (q *Queries) SetName(ctx context.Context, name string) error {
	return q.setName(ctx, name)
}
