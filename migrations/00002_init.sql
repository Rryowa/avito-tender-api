-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
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

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd