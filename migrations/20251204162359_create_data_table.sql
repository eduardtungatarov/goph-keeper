-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    data BYTEA NOT NULL

    INDEX idx_data_user_id (user_id),
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data;
-- +goose StatementEnd
