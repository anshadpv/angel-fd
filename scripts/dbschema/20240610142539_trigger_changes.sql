-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_insert_portfolio_from_webhook_events ON webhook_events;
create trigger trg_insert_portfolio_from_webhook_events after
insert on public.webhook_events for each row execute function insert_or_update_portfolio_from_webhook_events();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_insert_portfolio_from_webhook_events ON webhook_events;
-- +goose StatementEnd