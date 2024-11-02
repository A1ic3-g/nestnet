-- name: getName :one
SELECT name FROM LocalUser;

-- name: setName :exec
UPDATE LocalUser SET name=$1;

-- name: getPubKey :one
SELECT (pubX, pubY) FROM LocalUser;

-- name: getPrivKey :one
SELECT privD FROM LocalUser;