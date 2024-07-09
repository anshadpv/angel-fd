-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.insert_portfolio_from_webhook_events()
 RETURNS trigger LANGUAGE plpgsql
AS $function$
BEGIN
    IF NEW.client_code IS NOT NULL AND NEW.event_type = 'TD_BOOKED' THEN
        -- Check if the combination of client_code and vendor already exists in the portfolio table
        IF NOT EXISTS (
            SELECT 1 FROM portfolio
            WHERE client_code = NEW.client_code AND provider = NEW.vendor
        ) THEN
            -- Insert new record into portfolio table
            INSERT INTO portfolio (client_code, provider, invested_value, current_value, total_active_deposits, interest_earned, returns_value, returns_percentage, created_by) 
            VALUES (NEW.client_code, NEW.vendor, NEW.amount, NEW.amount, 1, 0, 0, 0, 'webhook_events');
        END IF;
    END IF;
    RETURN NULL; -- Trigger has completed successfully
END;
$function$
;

CREATE OR REPLACE FUNCTION public.insert_pending_journey()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    event_types TEXT[] := ARRAY['PAYMENT_FAILURE', 'AADHAAR_FAILED', 'VKYC_INITIATED', 'VKYC_FAILURE', 'VKYC_REQUIRED', 'VKYC_RETRY_REQUIRED', 'DIGILOCKER_FAILED', 'PAN_FAILURE'];
BEGIN
    -- Check if the new inserted row meets the conditions
    IF NEW.client_code IS NOT NULL AND NEW.event_type = ANY(event_types) THEN
        -- Check if the combination of client_code and vendor already exists in pending_journey
        IF NOT EXISTS (
            SELECT 1 FROM pending_journey 
            WHERE client_code = NEW.client_code AND provider = NEW.vendor
        ) THEN
            -- Insert the new row into pending_journey
            INSERT INTO pending_journey (client_code, provider, pending, payment_pending, kyc_pending, created_by)
            VALUES (NEW.client_code, NEW.vendor, FALSE, FALSE, FALSE, 'webhook_events');
        END IF;
    END IF;
    RETURN NEW;
END;
$function$
;

create trigger trg_insert_portfolio_from_webhook_events after
insert on public.webhook_events for each row execute function insert_portfolio_from_webhook_events();

create trigger insert_pending_journey_trigger after
insert on public.webhook_events for each row execute function insert_pending_journey();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_insert_portfolio_from_webhook_events ON webhook_events;
DROP TRIGGER IF EXISTS insert_pending_journey_trigger ON webhook_events;

DROP FUNCTION IF EXISTS public.insert_portfolio_from_webhook_events();
DROP FUNCTION IF EXISTS public.insert_pending_journey();
-- +goose StatementEnd
