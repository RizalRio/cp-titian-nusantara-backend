CREATE TABLE media_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_type VARCHAR(100) NOT NULL, -- Contoh: "Post", "Page", "Project"
    model_id UUID NOT NULL,           -- ID dari tabel induk
    media_type VARCHAR(50) NOT NULL,  -- Contoh: "thumbnail", "gallery"
    file_url TEXT NOT NULL,
    caption TEXT,
    "order" INTEGER DEFAULT 0,        -- Untuk urutan gambar jika ada banyak
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexing untuk mempercepat pencarian (Karena ini tabel polimorfik, index sangat penting!)
CREATE INDEX idx_media_assets_model ON media_assets(model_type, model_id);