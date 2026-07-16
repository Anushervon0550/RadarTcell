DROP INDEX IF EXISTS idx_technologies_deleted_at;

ALTER TABLE technologies
    DROP COLUMN IF EXISTS deleted_at;

