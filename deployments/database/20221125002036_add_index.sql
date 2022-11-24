-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX short_url_uniq
    ON shorten.url (short_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX shorten.short_url_uniq;
-- +goose StatementEnd
