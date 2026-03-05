ALTER TABLE metrics_definitions
    ADD COLUMN IF NOT EXISTS field_key TEXT;

-- Заполняем для уже существующих seed-метрик
UPDATE metrics_definitions
SET field_key = CASE lower(name)
                    WHEN 'trl' THEN 'readiness_level'
                    WHEN 'readiness level' THEN 'readiness_level'
                    WHEN 'readiness_level' THEN 'readiness_level'
                    WHEN 'custom metric 01' THEN 'custom_metric_1'
                    WHEN 'custom metric 02' THEN 'custom_metric_2'
                    WHEN 'custom metric 03' THEN 'custom_metric_3'
                    WHEN 'custom metric 04' THEN 'custom_metric_4'
                    ELSE field_key
    END
WHERE field_key IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'metrics_definitions_field_key_check'
    ) THEN
        ALTER TABLE metrics_definitions
            ADD CONSTRAINT metrics_definitions_field_key_check
                CHECK (
                    field_key IS NULL OR field_key IN (
                        'readiness_level',
                        'list_index',
                        'custom_metric_1',
                        'custom_metric_2',
                        'custom_metric_3',
                        'custom_metric_4'
                    )
                );
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS metrics_definitions_field_key_uidx
    ON metrics_definitions(field_key)
    WHERE field_key IS NOT NULL;