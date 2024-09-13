-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE TYPE organization_type AS ENUM (
--     'IE',
--     'LLC',
--     'JSC'
-- );
CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);
CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID UNIQUE REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type service_type NOT NULL,
    status tender_status DEFAULT 'Created',
    version INT DEFAULT 1,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id, version)
);

CREATE TABLE tender_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type service_type NOT NULL,
    status tender_status NOT NULL,
    version INT NOT NULL,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (tender_id, version)
);

-- Bid Tables
CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);
CREATE TYPE bid_decision AS ENUM (
    '',
    'Approved',
    'Rejected'
);
CREATE TYPE author_type_enum AS ENUM (
    'Organization',
    'User'
);

CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    feedback TEXT,
    status bid_status DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    decision bid_decision DEFAULT '',
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    author_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    author_type author_type_enum NOT NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id, version)
);

CREATE TABLE bid_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    feedback TEXT,
    status bid_status DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    decision bid_decision DEFAULT '',
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    author_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    author_type author_type_enum NOT NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id, version)
);

CREATE TABLE review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    author_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Create Tender functions and triggers
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
    current_version INT;
BEGIN
    SELECT * INTO tender_record
    FROM tender_history
    WHERE tender_history.tender_id = tenderId
      AND tender_history.version = rollback_version;

    SELECT t.version INTO current_version
    FROM tender t
    WHERE t.id = tenderId;

    UPDATE tender
    SET
        name = tender_record.name,
        description = tender_record.description,
        service_type = tender_record.service_type,
        status = tender_record.status,
        organization_id = tender_record.organization_id,
        creator_username = username
--         updated_at = CURRENT_TIMESTAMP,
--         version = current_version + 1
    WHERE tender.id = tenderId;

    RETURN QUERY
        SELECT tender.id, tender.name, tender.description, tender.service_type, tender.status, tender.version, tender.organization_id, tender.creator_username, tender.created_at, tender.updated_at
        FROM tender
        WHERE tender.id = tenderId;
END;
$$ LANGUAGE plpgsql;


-- Create Bid functions and triggers
CREATE OR REPLACE FUNCTION update_bid_metadata()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    NEW.version := NEW.version + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bid_metadata_trigger
    BEFORE UPDATE ON bid
    FOR EACH ROW
EXECUTE FUNCTION update_bid_metadata();



CREATE OR REPLACE FUNCTION save_bid_to_history()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO bid_history (bid_id, name, description, feedback, status, tender_id, organization_id, decision, author_id, author_type, version, created_at, updated_at)
    VALUES (NEW.id, NEW.name, NEW.description, NEW.feedback, NEW.status, NEW.tender_id, NEW.organization_id, NEW.decision, NEW.author_id, NEW.author_type, NEW.version, NEW.created_at, NEW.updated_at);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bid_insert_update_trigger
    AFTER INSERT OR UPDATE ON bid
    FOR EACH ROW
EXECUTE FUNCTION save_bid_to_history();

CREATE OR REPLACE FUNCTION rollback_bid_version(bidId UUID, rollback_version INT, username VARCHAR)
    RETURNS TABLE (id UUID, name VARCHAR(100), description TEXT, feedback TEXT, status bid_status, tender_id UUID, organization_id UUID, decision bid_decision, author_id UUID, author_type author_type_enum, version INT, created_at TIMESTAMP, updated_at TIMESTAMP) AS $$
DECLARE
    bid_record bid_history%ROWTYPE;
BEGIN
    SELECT * INTO bid_record
    FROM bid_history
    WHERE bid_history.bid_id = bidId AND bid_history.version = rollback_version;

    UPDATE bid
    SET
        name = bid_record.name,
        description = bid_record.description,
        feedback = bid_record.feedback,
        status = bid_record.status,
        tender_id = bid_record.tender_id,
        organization_id = bid_record.organization_id,
        decision = bid_record.decision,
        author_id = bid_record.author_id,
        author_username = username,
        author_type = bid_record.author_type,
        updated_at = CURRENT_TIMESTAMP,
        version = bid_record.version + 1
    WHERE bid.id = bidId;

    RETURN QUERY
        SELECT bid.id, bid.name, bid.description, bid.feedback, bid.status, bid.tender_id, bid.organization_id, bid.decision, bid.author_id, bid.author_type, bid.version, bid.created_at, bid.updated_at
        FROM bid
        WHERE bid.id = bidId;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_author_username()
    RETURNS TRIGGER AS $$
BEGIN
    SELECT username INTO NEW.author_username
    FROM employee
    WHERE id = NEW.author_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_author_username_trigger
    BEFORE INSERT ON bid
    FOR EACH ROW
EXECUTE FUNCTION set_author_username();


CREATE OR REPLACE FUNCTION set_organization_id()
    RETURNS TRIGGER AS $$
BEGIN
    SELECT organization_id INTO NEW.organization_id
    FROM tender
    WHERE id = NEW.tender_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_organization_id_trigger
    BEFORE INSERT ON bid
    FOR EACH ROW
EXECUTE FUNCTION set_organization_id();


-- Prepare data
INSERT INTO employee (username, first_name, last_name)
VALUES
    ('johndoe', 'John', 'Doe'),
    ('pepe', 'Pe', 'Pe'),
    ('robpike', 'Rob', 'Pike');

INSERT INTO organization (name, description, type)
VALUES
    ('Avito', 'A leading technology company.', 'LLC'),
    ('Ozon', 'A second technology company.', 'IE'),
    ('Yandex', 'A third technology company.', 'JSC');

INSERT INTO organization_responsible (organization_id, user_id)
VALUES
    ((SELECT id FROM organization WHERE name = 'Avito'), (SELECT id FROM employee WHERE username = 'johndoe')),
    ((SELECT id FROM organization WHERE name = 'Ozon'), (SELECT id FROM employee WHERE username = 'pepe')),
    ((SELECT id FROM organization WHERE name = 'Yandex'), (SELECT id FROM employee WHERE username = 'robpike'));
-- +goose StatementEnd