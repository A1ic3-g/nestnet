package generated

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"golang.org/x/net/context"
	"math/big"
)

func (q *Queries) GetName(ctx context.Context) (string, error) {
	return q.getName(ctx)
}

func (q *Queries) SetName(ctx context.Context, name string) error {
	return q.setName(ctx, name)
}

func (q *Queries) GetPrivKey(ctx context.Context) (big.Int, error) {
	dStr, err := q.getPrivKey(ctx)

	d := new(big.Int)
	d.SetString(string(dStr), 16)
	return *d, err
}

func (q *Queries) GetPubKey(ctx context.Context) (ecdsa.PublicKey, error) {
	keyInterface, err := q.getPubKey(ctx)
	keyValues := keyInterface.([2]big.Int)

	key := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &keyValues[0],
		Y:     &keyValues[1],
	}

	return key, err
}
