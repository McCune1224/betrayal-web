-- name: GetRoleByID :one
SELECT id, name, alignment, description, perks_json, starting_items_json
FROM roles
WHERE id = $1;

-- name: ListRoles :many
SELECT id, name, alignment, description, perks_json, starting_items_json
FROM roles
ORDER BY name;
