CREATE TABLE IF NOT EXISTS technology_metric_values (
    technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    metric_id UUID NOT NULL REFERENCES metrics_definitions(id) ON DELETE CASCADE,
    value DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (technology_id, metric_id)
);

-- Разрешаем динамические field_key (snake_case), чтобы не ограничивать набор метрик 1..4.
ALTER TABLE metrics_definitions
    DROP CONSTRAINT IF EXISTS metrics_definitions_field_key_check;

ALTER TABLE metrics_definitions
    ADD CONSTRAINT metrics_definitions_field_key_check
        CHECK (
            field_key IS NULL OR field_key ~ '^[a-z][a-z0-9_]{1,62}$'
        );

