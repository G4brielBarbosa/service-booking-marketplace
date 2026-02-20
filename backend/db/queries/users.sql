-- name: FindUserByTelegramID :one
SELECT * FROM user_profiles WHERE telegram_user_id = $1;

-- name: CreateUser :one
INSERT INTO user_profiles (user_id, telegram_user_id, primary_chat_id, timezone, locale, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUser :exec
UPDATE user_profiles
SET timezone = $2, locale = $3, updated_at = now()
WHERE user_id = $1;
