-- +goose Up
-- +goose StatementBegin
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
            INSERT INTO pending_journey (client_code, provider, pending, payment_pending, kyc_pending, created_by, updated_by)
            VALUES (NEW.client_code, NEW.vendor, FALSE, FALSE, FALSE, 'webhook_events', 'webhook_events');
        END IF;
    END IF;
    RETURN NEW;
END;
$function$
;
DROP TRIGGER IF EXISTS insert_pending_journey_trigger ON webhook_events;
create trigger insert_pending_journey_trigger after
insert on public.webhook_events for each row execute function insert_pending_journey();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
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
DROP TRIGGER IF EXISTS insert_pending_journey_trigger ON webhook_events;
create trigger insert_pending_journey_trigger after
insert on public.webhook_events for each row execute function insert_pending_journey();    
-- +goose StatementEnd
