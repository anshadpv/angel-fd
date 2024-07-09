-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION insert_or_update_portfolio_from_webhook_events()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.client_code IS NOT NULL AND NEW.event_type = 'TD_BOOKED' THEN
        -- Check if the combination of client_code and vendor already exists in the portfolio table
        IF EXISTS (
            SELECT 1 FROM portfolio
            WHERE client_code = NEW.client_code AND provider = NEW.vendor
        ) THEN
            -- Update the existing record in the portfolio table
            UPDATE portfolio
            SET
                invested_value = invested_value + NEW.amount,
                current_value = current_value + NEW.amount,
                total_active_deposits = total_active_deposits + 1,
                updated_at = CURRENT_TIMESTAMP,
                updated_by = 'webhook_event'
            WHERE
                client_code = NEW.client_code AND provider = NEW.vendor;
        ELSE
            -- Insert new record into portfolio table
            INSERT INTO portfolio (
                client_code, provider, invested_value, current_value, total_active_deposits, interest_earned, returns_value, returns_percentage, created_by, updated_by) 
                VALUES (NEW.client_code, NEW.vendor, NEW.amount, NEW.amount, 1, 0, 0, 0, 'webhook_event', 'webhook_event');
        END IF;
    END IF;
    RETURN NULL; -- Trigger has completed successfully
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_insert_portfolio_from_webhook_events ON webhook_events;

create trigger trg_insert_portfolio_from_webhook_events after
insert on public.webhook_events for each row execute function insert_portfolio_from_webhook_events();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION insert_or_update_portfolio_from_webhook_events()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.client_code IS NOT NULL AND NEW.event_type = 'TD_BOOKED' THEN
        -- Check if the combination of client_code and vendor already exists in the portfolio table
        IF EXISTS (
            SELECT 1 FROM portfolio
            WHERE client_code = NEW.client_code AND provider = NEW.vendor
        ) THEN
            -- Update the existing record in the portfolio table
            UPDATE portfolio
            SET
                invested_value = invested_value + NEW.amount,
                current_value = current_value + NEW.amount,
                total_active_deposits = total_active_deposits + 1,
                updated_at = CURRENT_TIMESTAMP,
                updated_by = 'webhook_event'
            WHERE
                client_code = NEW.client_code AND provider = NEW.vendor;
        ELSE
            -- Insert new record into portfolio table
            INSERT INTO portfolio (
                client_code, provider, invested_value, current_value, total_active_deposits, interest_earned, returns_value, returns_percentage, created_by) 
                VALUES (NEW.client_code, NEW.vendor, NEW.amount, NEW.amount, 1, 0, 0, 0, 'webhook_event');
        END IF;
    END IF;
    RETURN NULL; -- Trigger has completed successfully
END;
$$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS trg_insert_portfolio_from_webhook_events ON webhook_events;
create trigger trg_insert_portfolio_from_webhook_events after
insert on public.webhook_events for each row execute function insert_portfolio_from_webhook_events();
-- +goose StatementEnd
