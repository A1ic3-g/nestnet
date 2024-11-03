-- name: getPosts :many
SELECT * From Posts;

-- name: addPost :exec
INSERT INTO Posts (id, title, body, imgmd5) VALUES ($1, $2, $3, $4);