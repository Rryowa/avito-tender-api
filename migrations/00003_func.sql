-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_tender_metadata()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at := NOW();
        NEW.version := NEW.version + 1;
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

CREATE TRIGGER tender_metadata_trigger
    BEFORE UPDATE ON tender
    FOR EACH ROW
    EXECUTE FUNCTION update_tender_metadata();


CREATE OR REPLACE FUNCTION save_tender_to_history()
    RETURNS TRIGGER AS $$
    BEGIN
        INSERT INTO tender_history (tender_id, name, description, service_type, status, version, organization_id, creator_username, created_at, updated_at)
        VALUES (NEW.id, NEW.name, NEW.description, NEW.service_type, NEW.status, NEW.version, NEW.organization_id, NEW.creator_username, NEW.created_at, NEW.updated_at);
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

CREATE TRIGGER tender_insert_update_trigger
    AFTER INSERT OR UPDATE ON tender
    FOR EACH ROW
    EXECUTE FUNCTION save_tender_to_history();


CREATE OR REPLACE FUNCTION rollback_tender_version(tenderId UUID, rollback_version INT, username VARCHAR)
    RETURNS TABLE (id UUID, name VARCHAR, description TEXT, service_type service_type, status tender_status, version INT, organization_id UUID, creator_username VARCHAR, created_at TIMESTAMP, updated_at TIMESTAMP) AS $$
DECLARE
    tender_record tender_history%ROWTYPE;
BEGIN
    SELECT * INTO tender_record
    FROM tender_history
    WHERE tender_history.tender_id = tenderId AND tender_history.version = rollback_version;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Version % for tender % not found', rollback_version, tenderId;
    END IF;

    -- Update the current tender with the rollback version
    UPDATE tender
    SET
        name = tender_record.name,
        description = tender_record.description,
        service_type = tender_record.service_type,
        status = tender_record.status,
        organization_id = tender_record.organization_id,
        creator_username = username,
        updated_at = CURRENT_TIMESTAMP,
        version = tender_record.version + 1
    WHERE tender.id = tenderId;
    RETURN QUERY
        SELECT tender.id, tender.name, tender.description, tender.service_type, tender.status, tender.version, tender.organization_id, tender.creator_username, tender.created_at, tender.updated_at
        FROM tender
        WHERE tender.id = tenderId;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS tender_insert_update_trigger ON tender;
DROP FUNCTION IF EXISTS update_tender_updated_at();

DROP TRIGGER IF EXISTS tender_insert_update_trigger ON tender;
DROP FUNCTION IF EXISTS save_tender_to_history();

DROP FUNCTION IF EXISTS rollback_tender_version(tender_id UUID, rollback_version INT, username VARCHAR);
-- +goose StatementEnd