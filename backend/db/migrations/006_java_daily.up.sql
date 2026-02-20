-- Migration 006: Java Daily (PLAN-005)

CREATE TABLE IF NOT EXISTS java_practice_sessions (
    session_id          UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id             UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date          TEXT NOT NULL,
    duration_est_min    INT NOT NULL DEFAULT 0,
    objective_constraint TEXT NOT NULL DEFAULT '',
    evidence_short      TEXT NOT NULL DEFAULT '',
    status              TEXT NOT NULL DEFAULT 'partial',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_java_practice_user_date ON java_practice_sessions(user_id, local_date);
CREATE INDEX idx_java_practice_task ON java_practice_sessions(task_id);

CREATE TABLE IF NOT EXISTS java_retrievals (
    retrieval_id    UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id         UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date      TEXT NOT NULL,
    items_answered  INT NOT NULL DEFAULT 0,
    items_total     INT NOT NULL DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'low',
    targets         JSONB NOT NULL DEFAULT '[]',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_java_retrieval_user_date ON java_retrievals(user_id, local_date);
CREATE INDEX idx_java_retrieval_task ON java_retrievals(task_id);

CREATE TABLE IF NOT EXISTS java_learning_log (
    entry_id            UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id             UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date          TEXT NOT NULL,
    error_or_learning   TEXT NOT NULL,
    fix_or_note         TEXT,
    category            TEXT,
    recurring_count_14d INT NOT NULL DEFAULT 1,
    is_recurring        BOOLEAN NOT NULL DEFAULT false,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_java_learning_user_date ON java_learning_log(user_id, local_date);
CREATE INDEX idx_java_learning_user_label ON java_learning_log(user_id, error_or_learning);
