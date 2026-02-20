-- name: FindPrivacyPolicy :one
SELECT * FROM privacy_policies WHERE user_id = $1;

-- name: UpsertPrivacyPolicy :exec
INSERT INTO privacy_policies (user_id, opt_out_categories, retention_days, minimal_mode, updated_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (user_id) DO UPDATE SET
    opt_out_categories = EXCLUDED.opt_out_categories,
    retention_days = EXCLUDED.retention_days,
    minimal_mode = EXCLUDED.minimal_mode,
    updated_at = now();
