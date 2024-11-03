package generated

import "golang.org/x/net/context"

func (q *Queries) GetPeers(ctx context.Context) ([]Peer, error) {
	return q.getPeers(ctx)
}

func (q *Queries) AddPeer(ctx context.Context, peer *Peer) error {
	params := addPeerParams{
		ID:      peer.ID,
		Name:    peer.Name,
		Pubx:    peer.Pubx,
		Puby:    peer.Puby,
		Address: peer.Address,
	}
	return q.addPeer(ctx, params)
}
