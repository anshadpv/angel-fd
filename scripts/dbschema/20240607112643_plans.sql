-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plans (
	plan_id int4 NULL,
	fsi text NULL,
	plan_type text NULL,
	tenure_years int4 NULL,
	tenure_months int4 NULL,
	tenure_days int4 NULL,
	interest_rate numeric(10, 2) NULL,
	lockin_months int4 NULL,
	is_active bool NULL,
	is_insured bool NULL,
	is_mostbought bool NULL,
	senior_citizen_benefit numeric(10, 2) NULL,
	women_benefit numeric(10, 2) NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX plans_created ON plans (created_at);
CREATE INDEX plans_updated ON plans (updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX plans_created;
DROP INDEX plans_updated;
drop table plans;
-- +goose StatementEnd
