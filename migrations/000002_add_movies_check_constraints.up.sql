ALTER TABLE Movie ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE Movie ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND year(now()));
ALTER TABLE Movie ADD CONSTRAINT genres_length_check CHECK (length(genres, 1) BETWEEN 1 AND 5);
