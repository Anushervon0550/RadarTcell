DROP INDEX IF EXISTS metrics_definitions_field_key_uidx;

ALTER TABLE metrics_definitions
    DROP CONSTRAINT IF EXISTS metrics_definitions_field_key_check;

ALTER TABLE metrics_definitions
    DROP COLUMN IF EXISTS field_key;