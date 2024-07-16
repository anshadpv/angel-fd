-- +goose Up
-- +goose StatementBegin
ALTER TABLE portfolio_test
ADD COLUMN to_be_refreshed BOOLEAN NOT NULL DEFAULT false;
CREATE OR REPLACE FUNCTION public.insert_or_update_portfolio_test_from_wb_events_test()
 RETURNS trigger AS $$
DECLARE
    event_types TEXT[] := ARRAY['TD_BOOKED','BOOKING_SUCCESSFUL', 'SUCCESS'];
BEGIN
    IF NEW.client_code IS NOT NULL AND NEW.event_type = ANY(event_types) THEN
        -- Check if the combination of client_code and vendor already exists in the portfolio table
        IF EXISTS (
            SELECT 1 FROM portfolio_test
            WHERE client_code = NEW.client_code AND provider = NEW.vendor
        ) THEN
            -- Update the existing record in the portfolio table
            UPDATE portfolio
            SET
                invested_value = invested_value + NEW.amount,
                current_value = current_value + NEW.amount,
                total_active_deposits = total_active_deposits + 1,
                updated_at = CURRENT_TIMESTAMP,
                updated_by = 'wb_events_test',
                to_be_refreshed = true
            WHERE client_code = NEW.client_code AND provider = NEW.vendor;
        ELSE
            -- Insert new record into portfolio table
            INSERT INTO portfolio_test (
                client_code, provider, invested_value, current_value, total_active_deposits, interest_earned, returns_value, returns_percentage, created_by, updated_by,to_be_refreshed) 
                VALUES (NEW.client_code, NEW.vendor, NEW.amount, NEW.amount, 1, 0, 0, 0, 'wb_events_test', 'wb_events_test', true);
        END IF;
    END IF;
    RETURN NULL; -- Trigger has completed successfully

END;
$$ LANGUAGE plpgsql;
;

ALTER TABLE pending_journey_test
ADD COLUMN to_be_refreshed BOOLEAN NOT NULL DEFAULT false;
CREATE OR REPLACE FUNCTION public.insert_or_update_pending_journey_test_from_wb_events_test()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    event_types TEXT[] := ARRAY['TD_FAILED', 'BOOKING_UNSUCCESSFUL', 'KYC_PENDING'];
BEGIN
    -- Check if the new inserted row meets the conditions
    IF NEW.client_code IS NOT NULL AND NEW.event_type = ANY(event_types) THEN
        -- Check if the combination of client_code and vendor already exists in pending_journey
       IF EXISTS (
            SELECT 1 FROM pending_journey_test 
            WHERE client_code = NEW.client_code AND provider = NEW.vendor
        ) THEN
            -- Update the existing row in pending_journey
            UPDATE pending_journey
            SET updated_at = CURRENT_TIMESTAMP,
                updated_by = 'wb_events_test',
                to_be_refreshed = true
            WHERE client_code = NEW.client_code AND provider = NEW.vendor;
        ELSE
            -- Insert the new row into pending_journey
            INSERT INTO pending_journey_test (client_code, provider, pending, payment_pending, kyc_pending, created_by, updated_by, to_be_refreshed)
            VALUES (NEW.client_code, NEW.vendor, FALSE, FALSE, FALSE, 'wb_events_test', 'wb_events_test', true);
        END IF;
    END IF;
    RETURN NEW;
END;
$function$
;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE portfolio_test
DROP COLUMN to_be_refreshed;
ALTER TABLE pending_journey_test
DROP COLUMN to_be_refreshed;

DROP FUNCTION IF EXISTS public.insert_or_update_portfolio_test_from_wb_events_test() CASCADE;
DROP FUNCTION IF EXISTS public.insert_or_update_pending_journey_test_from_wb_events_test() CASCADE;

-- +goose StatementEnd
