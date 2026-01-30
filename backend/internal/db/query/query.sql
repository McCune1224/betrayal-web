-- Example SQLC query:
-- -- name: GetUserByID :one
select id, username, email, created_at
from users
where id = $1
;
