-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS faqs (
	tag varchar NOT NULL,
	faq json NOT NULL,
	created_ts timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_ts timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_active bool NULL DEFAULT true,
	CONSTRAINT faqs_pkey PRIMARY KEY (tag)
);
CREATE INDEX faqs_created ON faqs (created_ts);
CREATE INDEX faqs_updated ON faqs (updated_ts);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE faqs DROP CONSTRAINT faqs_pkey;
DROP INDEX faqs_created;
DROP INDEX faqs_updated;
drop table faqs;
-- +goose StatementEnd
