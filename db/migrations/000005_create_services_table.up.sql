CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    short_description TEXT,
    description TEXT,
    icon_name VARCHAR(100),
    is_flagship BOOLEAN DEFAULT false,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Digunakan untuk Soft Delete GORM
);

-- Membuat index untuk mempercepat query berdasarkan slug (pencarian halaman publik)
CREATE INDEX idx_services_slug ON services(slug);

-- Membuat index untuk mempercepat query data yang tidak dihapus (Soft Delete)
CREATE INDEX idx_services_deleted_at ON services(deleted_at);