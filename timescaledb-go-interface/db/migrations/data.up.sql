-- +migrate Up
CREATE extension IF NOT EXISTS "uuid-ossp";

-- +migrate Up
CREATE TABLE IF NOT EXISTS project_data(
   created_at timestamp NOT NULL DEFAULT NOW() PRIMARY KEY,
   project_id uuid NOT NULL,
   run_seq_no integer NOT NULL DEFAULT 0,
   data json
);

-- +migrate Up
CREATE INDEX ON project_data (created_at DESC, project_id);