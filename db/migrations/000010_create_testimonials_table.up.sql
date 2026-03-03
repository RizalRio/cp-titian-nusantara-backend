-- 1. Hapus kolom testimoni lama di tabel portfolios
ALTER TABLE portfolios 
DROP COLUMN IF EXISTS testimonial,
DROP COLUMN IF EXISTS testimonial_author;

-- 2. Buat tabel khusus testimonials
CREATE TABLE testimonials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portfolio_id UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    author_name VARCHAR(255) NOT NULL,
    author_role VARCHAR(255), -- Cth: "Kepala Desa", "Penerima Manfaat"
    content TEXT NOT NULL,
    avatar_url VARCHAR(255),  -- Opsional: URL foto pemberi testimoni
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_testimonials_portfolio_id ON testimonials(portfolio_id);