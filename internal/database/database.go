package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"nestnet/internal/database/generated"
)

// getQueries gets the queries struct used to query the database
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

// GetName gets the user's name
func GetName() string {
	name, err := getQueries().GetName(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return name
}

// SetName sets the user's name
func SetName(name string) {
	err := getQueries().SetName(context.Background(), name)
	if err != nil {
		log.Fatal(err)
	}
}

// GetPosts gets the user's posts
func GetPosts() []generated.Post {
	posts, err := getQueries().GetPosts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return posts
}

// GetPeers gets the user's peers
func GetPeers() []generated.Peer {
	peers, err := getQueries().GetPeers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return peers
}

// AddPost adds a post to the user's posts
func AddPost(post generated.Post) {
	err := getQueries().AddPost(context.Background(), post)
	if err != nil {
		log.Fatal(err)
	}
}

// AddPeer adds a peer to the user's peers
func AddPeer(peer generated.Peer) {
	err := getQueries().AddPeer(context.Background(), &peer)
	if err != nil {
		log.Fatal(err)
	}
}
