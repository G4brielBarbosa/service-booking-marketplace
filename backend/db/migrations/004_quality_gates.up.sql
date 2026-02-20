-- PLAN-003: Quality Gates & Evidence entities
-- Retention: C3 raw 7 days, C2 evidence 90 days, derived (gate_results, rubric_scores) 90 days

CREATE TABLE evidences (
    evidence_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    task_id        UUID NOT NULL REFERENCES planned_tasks(task_id) ON DELETE CASCADE,
    kind           TEXT NOT NULL CHECK (kind IN ('text_answer','rubric','audio','metadata')),
    sensitivity    TEXT NOT NULL CHECK (sensitivity IN ('C2','C3')),
    storage_policy TEXT NOT NULL CHECK (storage_policy IN ('kept_7d','kept_custom','discarded_after_processing')),
    content_ref    TEXT,
    summary        TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_evidences_task ON evidences (user_id, task_id);

CREATE TABLE gate_results (
    gate_result_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    task_id            UUID NOT NULL REFERENCES planned_tasks(task_id) ON DELETE CASCADE,
    gate_status        TEXT NOT NULL CHECK (gate_status IN ('satisfied','not_satisfied')),
    failure_reason_code TEXT NOT NULL DEFAULT '',
    reason_short       TEXT NOT NULL DEFAULT '',
    next_min_step      TEXT NOT NULL DEFAULT '',
    evidence_ids       UUID[] NOT NULL DEFAULT '{}',
    derived_metrics    JSONB NOT NULL DEFAULT '{}',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_gate_results_task ON gate_results (user_id, task_id);

CREATE TABLE rubric_scores (
    rubric_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    task_id    UUID NOT NULL REFERENCES planned_tasks(task_id) ON DELETE CASCADE,
    domain     TEXT NOT NULL,
    dimensions JSONB NOT NULL DEFAULT '{}',
    total      INT NOT NULL DEFAULT 0,
    status     TEXT NOT NULL CHECK (status IN ('complete','partial')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_rubric_scores_task ON rubric_scores (user_id, task_id);
