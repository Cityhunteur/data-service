BEGIN;

CREATE TABLE Data
(
    id        UUID PRIMARY KEY,
    title     VARCHAR(64),
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

END;