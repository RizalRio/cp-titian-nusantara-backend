-- Menggunakan tipe data JSONB agar Frontend mudah membaca array of string
ALTER TABLE services 
ADD COLUMN impact_points JSONB DEFAULT '[]'::jsonb,
ADD COLUMN cta_text VARCHAR(100),
ADD COLUMN cta_link VARCHAR(255);