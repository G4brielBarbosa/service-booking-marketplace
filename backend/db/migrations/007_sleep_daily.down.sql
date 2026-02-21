-- Migration 007 rollback: Sleep Daily

DROP TABLE IF EXISTS weekly_sleep_interventions;
DROP TABLE IF EXISTS sleep_routine_records;
DROP TABLE IF EXISTS sleep_diary_entries;
