-- name: GetUserFromRefreshToken :one
SELECT u.id, u.email, u.hashed_password, u.created_at, u.updated_at
FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1
  AND rt.revoked_at IS NULL
  AND rt.expires_at > NOW();
  