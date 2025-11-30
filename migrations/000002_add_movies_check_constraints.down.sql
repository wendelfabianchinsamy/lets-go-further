ALTER TABLE MOVIES
DROP CONSTRAINT movies_runtime_check IF EXISTS movies_runtime_check;

ALTER TABLE MOVIES
DROP CONSTRAINT movies_year_check IF EXISTS movies_year_check;

ALTER TABLE MOVIES
DROP CONSTRAINT movies_genres_length_check IF EXISTS movies_genres_length_check;