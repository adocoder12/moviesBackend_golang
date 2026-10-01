CREATE UNIQUE INDEX IF NOT EXISTS idx_movies_title_year
ON movies(title COLLATE NOCASE, year);
