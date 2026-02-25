-- Hapus dari tabel pivot / yang paling banyak memiliki foreign key terlebih dahulu
DROP TABLE IF EXISTS post_tags;

-- Hapus index (opsional di PostgreSQL karena akan terhapus bersama tabel, tapi praktik yang baik)
DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_category_id;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_deleted_at;

-- Hapus tabel utama
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;