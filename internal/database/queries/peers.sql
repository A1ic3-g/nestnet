-- name: getPeers :many
SELECT * FROM Peers;

-- name: addPeer :exec
INSERT INTO Peers (id, name, pubX, pubY, address) VALUES ($1, $2, $3, $4, $5);