-- +goose Up
-- Minimal reference table backing the seeded "example" domain (see
-- internal/domain/example). Keep this migration file single-purpose per
-- service while pre-launch: add new tables here rather than a new file.
CREATE TABLE IF NOT EXISTS examples (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS examples;
