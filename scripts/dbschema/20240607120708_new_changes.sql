-- +goose Up
-- +goose StatementBegin
ALTER TABLE pending_journey
            ALTER COLUMN created_by SET NOT NULL,
            ALTER COLUMN updated_by SET NOT NULL;
ALTER TABLE pending_journey
            ADD COLUMN invalid_client bool,
            ADD COLUMN api_error varchar(100) DEFAULT '';
ALTER TABLE portfolio
            ADD COLUMN invalid_client bool,
            ADD COLUMN api_error varchar(100) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pending_journey
            ALTER COLUMN created_by DROP NOT NULL,
            ALTER COLUMN updated_by DROP NOT NULL;
ALTER TABLE pending_journey
            DROP COLUMN invalid_client,
            DROP COLUMN api_error;
ALTER TABLE portfolio
            DROP COLUMN invalid_client,
            DROP COLUMN api_error;
-- +goose StatementEnd
