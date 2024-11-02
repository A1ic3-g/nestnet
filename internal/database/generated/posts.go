package generated

import "golang.org/x/net/context"

func (q *Queries) GetPosts(ctx context.Context) ([]Post, error) {
	return q.getPosts(ctx)
}
