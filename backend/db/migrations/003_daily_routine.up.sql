-- PLAN-002: Daily routine entities (DailyState, DailyCheckIn, DailyPlan, PlannedTask)
-- Retention: 90 days (C1) managed by application/worker

CREATE TABLE daily_states (
    user_id     UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    local_date  DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, local_date)
);

CREATE TABLE daily_check_ins (
    check_in_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    local_date       DATE NOT NULL,
    time_available_min INT NOT NULL,
    energy_0_10      INT NOT NULL CHECK (energy_0_10 BETWEEN 0 AND 10),
    mood_stress_0_10 INT CHECK (mood_stress_0_10 BETWEEN 0 AND 10),
    constraints_text TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (user_id, local_date)
);

CREATE TABLE daily_plans (
    plan_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id              UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    local_date           DATE NOT NULL,
    plan_type            TEXT NOT NULL CHECK (plan_type IN ('A','B','C')),
    rationale            TEXT NOT NULL DEFAULT '',
    priority_task_id     UUID NOT NULL,
    complementary_ids    UUID[] NOT NULL DEFAULT '{}',
    foundation_task_id   UUID,
    version              INT NOT NULL DEFAULT 1,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_daily_plans_user_date ON daily_plans (user_id, local_date);

CREATE TABLE planned_tasks (
    task_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    local_date     DATE NOT NULL,
    title          TEXT NOT NULL,
    goal_domain    TEXT NOT NULL,
    estimated_min  INT NOT NULL DEFAULT 5,
    instructions   TEXT NOT NULL DEFAULT '',
    done_criteria  TEXT NOT NULL DEFAULT '',
    status         TEXT NOT NULL DEFAULT 'planned'
                   CHECK (status IN ('planned','in_progress','completed','blocked','deferred','evidence_pending','attempt')),
    block_reason   TEXT,
    note           TEXT,
    gate_ref       TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_planned_tasks_user_date ON planned_tasks (user_id, local_date);
