
-- +migrate Up
CREATE TABLE IF NOT EXISTS bitlytest(
    id SERIAL8 PRIMARY KEY,
    small_url TEXT NOT NULL CHECK (small_url <> '') UNIQUE,
    origin_url TEXT NOT NULL CHECK (origin_url <> ''),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX ON bitlytest (small_url);

-- +migrate Down
DROP TABLE bitlytest;
