-- +goose Up
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS review;
DROP TABLE IF EXISTS bid_history;
DROP TABLE IF EXISTS bid;

DROP TABLE IF EXISTS tender_history;
DROP TABLE IF EXISTS tender;
DROP TABLE IF EXISTS organization_responsible;
DROP TABLE IF EXISTS organization;
DROP TABLE IF EXISTS employee;

DROP TRIGGER IF EXISTS tender_insert_update_trigger ON tender;
DROP FUNCTION IF EXISTS update_tender_updated_at();

DROP TRIGGER IF EXISTS tender_insert_update_trigger ON tender;
DROP FUNCTION IF EXISTS save_tender_to_history();
DROP FUNCTION IF EXISTS rollback_tender_version(tender_id UUID, rollback_version INT, username VARCHAR);

DROP TYPE IF EXISTS service_type;
DROP TYPE IF EXISTS tender_status;
DROP TYPE IF EXISTS organization_type;

DROP TRIGGER IF EXISTS set_author_username_trigger ON bid;
DROP FUNCTION IF EXISTS set_author_username();

DROP TRIGGER IF EXISTS bid_metadata_trigger ON bid;
DROP FUNCTION IF EXISTS update_bid_metadata();

DROP TRIGGER IF EXISTS bid_insert_update_trigger ON bid;
DROP FUNCTION IF EXISTS save_bid_to_history();
DROP FUNCTION IF EXISTS rollback_bid_version(bidId UUID, rollback_version INT, username VARCHAR);

DROP TRIGGER IF EXISTS set_organization_id ON bid;
DROP FUNCTION IF EXISTS set_organization_id();

DROP TYPE IF EXISTS bid_decision;
DROP TYPE IF EXISTS author_type_enum;
DROP TYPE IF EXISTS bid_status;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd