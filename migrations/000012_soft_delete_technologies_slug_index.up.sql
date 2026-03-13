ALTER TABLE technologies
    DROP CONSTRAINT IF EXISTS technologies_slug_key;

CREATE UNIQUE INDEX IF NOT EXISTS technologies_slug_active_uidx
    ON technologies(slug)
    WHERE deleted_at IS NULL;

