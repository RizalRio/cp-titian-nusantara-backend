ALTER TABLE portfolios
ALTER COLUMN locations TYPE VARCHAR(255);

ALTER TABLE portfolios
RENAME COLUMN locations TO location;