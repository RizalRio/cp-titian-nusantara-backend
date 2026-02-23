CREATE TABLE site_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- (Opsional tapi direkomendasikan) Masukkan data default (Seeding)
INSERT INTO site_settings (key, value) VALUES 
('contact_email', 'halo@titiannusantara.com'),
('footer_manifesto', 'Titian Nusantara: Melangkah bersama untuk dampak yang lebih luas.'),
('social_instagram', 'https://instagram.com/titiannusantara');