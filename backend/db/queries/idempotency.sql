-- name: CheckIdempotency :one
SELECT * FROM idempotency_records WHERE key = $1 AND expires_at > now();

-- name: StoreIdempotency :exec
INSERT INTO idempotency_records (key, first_seen_at, result_ref, expires_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (key) DO NOTHING;

-- name: CleanExpiredIdempotency :exec
DELETE FROM idempotency_records WHERE expires_at <= now();
