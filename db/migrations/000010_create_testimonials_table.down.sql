DROP TABLE IF EXISTS testimonials;

ALTER TABLE portfolios 
ADD COLUMN testimonial TEXT,
ADD COLUMN testimonial_author VARCHAR(255);