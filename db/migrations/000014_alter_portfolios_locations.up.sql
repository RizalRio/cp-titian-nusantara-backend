ALTER TABLE portfolios
RENAME COLUMN location TO locations;

ALTER TABLE portfolios
ALTER COLUMN locations TYPE JSONB
USING NULL;

UPDATE portfolios
SET locations = NULL;