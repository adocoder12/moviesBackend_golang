-- +goose Up
CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    year INTEGER NOT NULL,
    director TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS movies;
