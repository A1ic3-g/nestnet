package generated

import "golang.org/x/net/context"

func (q *Queries) GetPeers(ctx context.Context) ([]Peer, error) {
	return q.getPeers(ctx)
}
