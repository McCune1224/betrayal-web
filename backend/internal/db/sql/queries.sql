-- name: GetRoleByID :one
SELECT id, name, alignment, description, perks_json, starting_items_json
FROM roles
WHERE id = $1;

-- name: ListRoles :many
SELECT id, name, alignment, description, perks_json, starting_items_json
FROM roles
ORDER BY name;

-- name: CreateRoom :one
INSERT INTO rooms (code, host_id, phase)
VALUES ($1, $2, $3)
RETURNING code, host_id, phase, created_at;

-- name: GetRoomByCode :one
SELECT code, host_id, phase, created_at FROM rooms WHERE code = $1;

-- name: CreatePlayer :one
INSERT INTO players (id, name, room_code, role_id, is_alive)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, room_code, role_id, is_alive, joined_at;

-- name: GetPlayersInRoom :many
SELECT id, name, room_code, role_id, is_alive, joined_at
FROM players WHERE room_code = $1;

-- name: CreateAction :one
INSERT INTO actions (id, player_id, type, target_id, room_code, phase)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, player_id, type, target_id, room_code, phase, timestamp;

-- name: GetActionsForRoomPhase :many
SELECT id, player_id, type, target_id, room_code, phase, timestamp
FROM actions WHERE room_code = $1 AND phase = $2;
