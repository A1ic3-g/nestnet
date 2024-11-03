-- name: getName :one
SELECT name FROM LocalUser;

-- name: setName :exec
UPDATE LocalUser SET name=$1;