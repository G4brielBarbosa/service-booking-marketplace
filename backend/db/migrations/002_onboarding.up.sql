-- PLAN-001: Onboarding entities

CREATE TABLE onboarding_sessions (
    session_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    status          TEXT NOT NULL DEFAULT 'new'
                    CHECK (status IN ('new','in_progress','minimum_completed','completed')),
    current_step_id TEXT NOT NULL DEFAULT 'welcome',
    answers         JSONB NOT NULL DEFAULT '[]',
    pending_items   JSONB NOT NULL DEFAULT '[]',
    started_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_interaction_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ
);

CREATE TABLE active_goal_cycles (
    cycle_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL UNIQUE REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    active_goals  JSONB NOT NULL DEFAULT '[]',
    paused_goals  JSONB NOT NULL DEFAULT '[]',
    started_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE baseline_snapshots (
    baseline_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    domain        TEXT NOT NULL,
    data          JSONB NOT NULL DEFAULT '{}',
    completeness  TEXT NOT NULL DEFAULT 'minimum'
                  CHECK (completeness IN ('minimum','partial','complete')),
    captured_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(user_id, domain)
);

CREATE TABLE minimum_viable_dailies (
    mvd_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL UNIQUE REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    items       JSONB NOT NULL DEFAULT '[]',
    when_to_use TEXT NOT NULL DEFAULT '',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
