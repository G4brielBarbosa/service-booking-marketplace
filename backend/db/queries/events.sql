-- name: AppendEvent :exec
INSERT INTO domain_event_log (event_id, user_id, timestamp, local_date, week_id, event_type, payload_min, sensitivity)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
