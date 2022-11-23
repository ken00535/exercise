-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA shorten;
CREATE TABLE shorten.url (
    id SERIAL PRIMARY KEY,
    url VARCHAR(256) NOT NULL,
	short_url VARCHAR(64) NOT NULL,
	expire_at timestamptz NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA shorten CASCADE;
-- +goose StatementEnd
