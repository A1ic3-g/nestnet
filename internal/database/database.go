package database

import (
	"context"
	"crypto/ecdsa"
	"github.com/jackc/pgx/v5"
	"log"
	"nestnet/internal/database/generated"
)

func getQueries() *generated.Queries {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=nest password=net dbname=nestnetdatabase host=db port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	/*
		defer func(conn *pgx.Conn, ctx context.Context) {
			err := conn.Close(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}(conn, ctx)
	*/

	queries := generated.New(conn)
	return queries
}

func GetName() string {
	name, err := getQueries().GetName(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return name
}

func SetName(name string) {
	err := getQueries().SetName(context.Background(), name)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPubKey() ecdsa.PublicKey {
	key, err := getQueries().GetPubKey(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return key
}

func GetPrivKey() ecdsa.PrivateKey {
	d, err := getQueries().GetPrivKey(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return ecdsa.PrivateKey{
		PublicKey: GetPubKey(),
		D:         &d,
	}
}

func GetPosts() []generated.Post {
	posts, err := getQueries().GetPosts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return posts
}

func GetPeers() []generated.Peer {
	peers, err := getQueries().GetPeers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return peers
}

func AddPost(post generated.Post) {
	err := getQueries().AddPost(context.Background(), post)
	if err != nil {
		log.Fatal(err)
	}
}

func AddPeer(peer generated.Peer) {
	err := getQueries().AddPeer(context.Background(), &peer)
	if err != nil {
		log.Fatal(err)
	}
}
