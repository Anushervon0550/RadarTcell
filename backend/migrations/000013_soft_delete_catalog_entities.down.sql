DROP INDEX IF EXISTS idx_metrics_definitions_deleted_at;
DROP INDEX IF EXISTS idx_organizations_deleted_at;
DROP INDEX IF EXISTS idx_tags_deleted_at;
DROP INDEX IF EXISTS idx_sdgs_deleted_at;
DROP INDEX IF EXISTS idx_trends_deleted_at;

ALTER TABLE metrics_definitions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE organizations DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE tags DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE sdgs DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE trends DROP COLUMN IF EXISTS deleted_at;

