CREATE SCHEMA IF NOT EXISTS %s;

CREATE TABLE IF NOT EXISTS schema_version (
    version       int NOT NULL PRIMARY KEY,
    created_at    timestamp without time zone NOT NULL DEFAULT current_timestamp
);

