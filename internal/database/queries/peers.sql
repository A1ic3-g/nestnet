-- name: getPeers :many
SELECT * FROM Peers;

-- name: addPeer :exec
INSERT INTO Peers (id, name, address) VALUES ($1, $2, $3);