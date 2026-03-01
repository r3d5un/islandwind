CREATE INDEX IF NOT EXISTS idx_cache_general_id ON cache.general USING HASH (id);
CREATE INDEX IF NOT EXISTS idx_cache_general_expires_at ON cache.general (expires_at);
