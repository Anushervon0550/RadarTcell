ALTER TABLE technologies
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_technologies_deleted_at
    ON technologies(deleted_at);

