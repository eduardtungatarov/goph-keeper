-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    data BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_data_user_id ON data(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_data_user_id;
DROP TABLE IF EXISTS data;
-- +goose StatementEnd