-- +goose Up
-- +goose StatementBegin
UPDATE public.portfolio_test
SET invalid_client = false
WHERE invalid_client IS NULL;

UPDATE public.pending_journey_test
SET invalid_client = false
WHERE invalid_client IS NULL;
ALTER TABLE pending_journey_test
            ALTER COLUMN api_error SET DEFAULT NULL;

ALTER TABLE portfolio_test
            ALTER COLUMN api_error SET DEFAULT NULL;

ALTER TABLE public.portfolio_test
            ALTER COLUMN invalid_client SET NOT NULL,
            ALTER COLUMN invalid_client SET DEFAULT false;

ALTER TABLE public.pending_journey_test
            ALTER COLUMN invalid_client SET NOT NULL,
            ALTER COLUMN invalid_client SET DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pending_journey_test
            ALTER COLUMN api_error DROP DEFAULT;

ALTER TABLE portfolio_test
            ALTER COLUMN api_error DROP DEFAULT;

ALTER TABLE public.portfolio_test
            ALTER COLUMN invalid_client DROP NOT NULL,
            ALTER COLUMN invalid_client DROP DEFAULT;

ALTER TABLE public.pending_journey_test
            ALTER COLUMN invalid_client DROP NOT NULL,
            ALTER COLUMN invalid_client DROP DEFAULT;
-- +goose StatementEnd
