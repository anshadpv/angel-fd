-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_insert_portfolio_test_from_wb_events_test ON wb_events_test;
create trigger trg_insert_portfolio_test_from_wb_events_test after
insert on public.wb_events_test for each row execute function public.insert_or_update_portfolio_test_from_wb_events_test();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_insert_portfolio_test_from_wb_events_test ON wb_events_test;
-- +goose StatementEnd
