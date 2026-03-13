DROP TABLE IF EXISTS technology_metric_values;

ALTER TABLE metrics_definitions
    DROP CONSTRAINT IF EXISTS metrics_definitions_field_key_check;

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

