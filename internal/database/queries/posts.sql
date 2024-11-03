-- name: getPosts :many
SELECT * From Posts;

-- name: addPost :exec
INSERT INTO Posts (id, title, body, imgmd5, imgname) VALUES ($1, $2, $3, $4, $5);