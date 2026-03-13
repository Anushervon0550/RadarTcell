CREATE INDEX IF NOT EXISTS idx_technologies_slug_trgm
    ON technologies USING gin (slug gin_trgm_ops);

