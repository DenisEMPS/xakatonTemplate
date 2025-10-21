-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS template (
    id bigint PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS template;
-- +goose StatementEnd
