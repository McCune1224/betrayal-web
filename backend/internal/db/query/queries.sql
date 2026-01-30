-- name: CreateRoom :one
INSERT INTO rooms (code, host_id, phase)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoomByCode :one
SELECT * FROM rooms
WHERE code = $1;

-- name: GetRoomByID :one
SELECT * FROM rooms
WHERE id = $1;

-- name: UpdateRoomPhase :exec
UPDATE rooms
SET phase = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;

-- name: CreatePlayer :one
INSERT INTO players (room_id, name, role_id, is_alive)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPlayerByID :one
SELECT * FROM players
WHERE id = $1;

-- name: GetPlayersByRoom :many
SELECT * FROM players
WHERE room_id = $1
ORDER BY joined_at;

-- name: UpdatePlayerAlive :exec
UPDATE players
SET is_alive = $2
WHERE id = $1;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1;

-- name: GetAllRoles :many
SELECT * FROM roles
ORDER BY name;

-- name: CreateAction :one
INSERT INTO actions (room_id, player_id, action_type, target_id, phase)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetActionsByRoom :many
SELECT * FROM actions
WHERE room_id = $1
ORDER BY created_at;

-- name: GetActionsByRoomAndPhase :many
SELECT * FROM actions
WHERE room_id = $1 AND phase = $2
ORDER BY created_at;

-- name: DeleteAction :exec
DELETE FROM actions
WHERE id = $1;
