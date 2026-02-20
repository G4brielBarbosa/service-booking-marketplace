-- Migration 005: English Daily (PLAN-004)

CREATE TABLE IF NOT EXISTS english_input_sessions (
    session_id    UUID PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id       UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date    TEXT NOT NULL,
    duration_est_min INT NOT NULL DEFAULT 0,
    content_descriptor TEXT NOT NULL DEFAULT '',
    comprehension_answers JSONB NOT NULL DEFAULT '[]',
    status        TEXT NOT NULL DEFAULT 'partial',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_english_input_user_date ON english_input_sessions(user_id, local_date);
CREATE INDEX idx_english_input_task ON english_input_sessions(task_id);

CREATE TABLE IF NOT EXISTS english_retrievals (
    retrieval_id  UUID PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES user_profiles(user_id),
    task_id       UUID NOT NULL REFERENCES planned_tasks(task_id),
    local_date    TEXT NOT NULL,
    items_answered INT NOT NULL DEFAULT 0,
    items_total   INT NOT NULL DEFAULT 0,
    status        TEXT NOT NULL DEFAULT 'low',
    targets       JSONB NOT NULL DEFAULT '[]',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_english_retrieval_user_date ON english_retrievals(user_id, local_date);
CREATE INDEX idx_english_retrieval_task ON english_retrievals(task_id);

CREATE TABLE IF NOT EXISTS english_error_log (
    error_id      UUID PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES user_profiles(user_id),
    local_date    TEXT NOT NULL,
    label         TEXT NOT NULL,
    note_short    TEXT,
    recurring_count_14d INT NOT NULL DEFAULT 1,
    is_recurring  BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_english_error_user_date ON english_error_log(user_id, local_date);
CREATE INDEX idx_english_error_user_label ON english_error_log(user_id, label);
