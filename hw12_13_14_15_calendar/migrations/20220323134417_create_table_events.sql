-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id            INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id       INT UNSIGNED     NOT NULL,
    title         VARCHAR(255)     NOT NULL,
    event_time    DATETIME         NOT NULL,
    duration_min  TINYINT UNSIGNED NOT NULL,
    notice_before TINYINT   DEFAULT NULL,
    description   TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = INNODB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
