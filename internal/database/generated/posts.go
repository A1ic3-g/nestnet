package generated

import (
	"golang.org/x/net/context"
)

func (q *Queries) GetPosts(ctx context.Context) ([]Post, error) {
	return q.getPosts(ctx)
}

func (q *Queries) AddPost(ctx context.Context, post Post) error {
	params := addPostParams{
		ID:      post.ID,
		Title:   post.Title,
		Body:    post.Body,
		Imgmd5:  post.Imgmd5,
		Imgname: post.Imgname,
	}
	return q.addPost(ctx, params)
}
