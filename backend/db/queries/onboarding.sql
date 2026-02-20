-- name: FindOnboardingByUserID :one
SELECT * FROM onboarding_sessions WHERE user_id = $1;

-- name: CreateOnboarding :exec
INSERT INTO onboarding_sessions (session_id, user_id, status, current_step_id, answers, pending_items, started_at, last_interaction_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateOnboarding :exec
UPDATE onboarding_sessions
SET status = $2, current_step_id = $3, answers = $4, pending_items = $5,
    last_interaction_at = $6, completed_at = $7
WHERE session_id = $1;

-- name: FindGoalCycleByUserID :one
SELECT * FROM active_goal_cycles WHERE user_id = $1;

-- name: CreateGoalCycle :exec
INSERT INTO active_goal_cycles (cycle_id, user_id, active_goals, paused_goals, started_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateGoalCycle :exec
UPDATE active_goal_cycles
SET active_goals = $2, paused_goals = $3, updated_at = now()
WHERE cycle_id = $1;

-- name: FindBaselineByUserDomain :one
SELECT * FROM baseline_snapshots WHERE user_id = $1 AND domain = $2;

-- name: FindBaselinesByUser :many
SELECT * FROM baseline_snapshots WHERE user_id = $1;

-- name: UpsertBaseline :exec
INSERT INTO baseline_snapshots (baseline_id, user_id, domain, data, completeness, captured_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id, domain) DO UPDATE SET
    data = EXCLUDED.data,
    completeness = EXCLUDED.completeness,
    updated_at = now();

-- name: FindMVDByUserID :one
SELECT * FROM minimum_viable_dailies WHERE user_id = $1;

-- name: UpsertMVD :exec
INSERT INTO minimum_viable_dailies (mvd_id, user_id, items, when_to_use, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE SET
    items = EXCLUDED.items,
    when_to_use = EXCLUDED.when_to_use,
    updated_at = now();
