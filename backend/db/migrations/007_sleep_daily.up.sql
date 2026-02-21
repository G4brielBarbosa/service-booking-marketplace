-- Migration 007: Sleep Daily (PLAN-006)

CREATE TABLE IF NOT EXISTS sleep_diary_entries (
    entry_id             UUID PRIMARY KEY,
    user_id              UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id              UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date           TEXT NOT NULL,
    slept_at             TEXT,
    woke_at              TEXT,
    quality_0_10         INT,
    morning_energy_0_10  INT,
    computed_duration_min INT,
    awakenings_note      TEXT,
    status               TEXT NOT NULL DEFAULT 'partial',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sleep_diary_user_date ON sleep_diary_entries(user_id, local_date);
CREATE INDEX idx_sleep_diary_task ON sleep_diary_entries(task_id);

CREATE TABLE IF NOT EXISTS sleep_routine_records (
    record_id   UUID PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id     UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date  TEXT NOT NULL,
    version     TEXT NOT NULL DEFAULT 'normal',
    steps_done  JSONB NOT NULL DEFAULT '[]',
    result      TEXT NOT NULL DEFAULT 'not_done',
    note_short  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sleep_routine_user_date ON sleep_routine_records(user_id, local_date);
CREATE INDEX idx_sleep_routine_task ON sleep_routine_records(task_id);

CREATE TABLE IF NOT EXISTS weekly_sleep_interventions (
    intervention_id      UUID PRIMARY KEY,
    user_id              UUID NOT NULL REFERENCES user_profiles(user_id),
    week_id              TEXT NOT NULL,
    description          TEXT NOT NULL,
    why_short            TEXT NOT NULL DEFAULT '',
    adherence_rule       TEXT NOT NULL DEFAULT '',
    status               TEXT NOT NULL DEFAULT 'proposed',
    adherence_count_done INT NOT NULL DEFAULT 0,
    closing_outcome      TEXT,
    closed_at            TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sleep_intervention_user_week ON weekly_sleep_interventions(user_id, week_id);
