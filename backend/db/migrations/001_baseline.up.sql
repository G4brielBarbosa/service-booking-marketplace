-- PLAN-000: Platform baseline schema
-- Entities: UserProfile, PrivacyPolicy, DomainEventLog, IdempotencyRecord

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE user_profiles (
    user_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_user_id BIGINT NOT NULL UNIQUE,
    primary_chat_id  BIGINT NOT NULL,
    timezone       TEXT NOT NULL DEFAULT 'America/Sao_Paulo',
    locale         TEXT NOT NULL DEFAULT 'pt-BR',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE privacy_policies (
    user_id            UUID PRIMARY KEY REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    opt_out_categories TEXT[] NOT NULL DEFAULT '{}',
    retention_days     JSONB NOT NULL DEFAULT '{"C1":90,"C2":90,"C3":7,"C4":365,"C5":0}',
    minimal_mode       BOOLEAN NOT NULL DEFAULT false,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE domain_event_log (
    event_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    timestamp      TIMESTAMPTZ NOT NULL DEFAULT now(),
    local_date     DATE NOT NULL,
    week_id        TEXT NOT NULL,
    event_type     TEXT NOT NULL,
    payload_min    JSONB NOT NULL DEFAULT '{}',
    sensitivity    TEXT NOT NULL DEFAULT 'C1',

    CONSTRAINT chk_sensitivity CHECK (sensitivity IN ('C1','C2','C3','C4','C5'))
);

CREATE INDEX idx_events_user_date ON domain_event_log (user_id, local_date);
CREATE INDEX idx_events_type ON domain_event_log (event_type);
CREATE INDEX idx_events_week ON domain_event_log (user_id, week_id);

CREATE TABLE idempotency_records (
    key           TEXT PRIMARY KEY,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    result_ref    TEXT NOT NULL DEFAULT '',
    expires_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_idempotency_expires ON idempotency_records (expires_at);
