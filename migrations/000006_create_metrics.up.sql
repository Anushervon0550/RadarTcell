CREATE TABLE IF NOT EXISTS metrics_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL UNIQUE,
    description TEXT,

    type TEXT NOT NULL CHECK (type IN ('distance', 'bubble', 'bar')),
    orderable BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

-- значения метрик по технологиям (расширяемая модель)
CREATE TABLE IF NOT EXISTS technology_metric_values (
    technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    metric_id      UUID NOT NULL REFERENCES metrics_definitions(id) ON DELETE CASCADE,
    value          NUMERIC NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (technology_id, metric_id)
    );

CREATE INDEX IF NOT EXISTS idx_metric_values_metric_id ON technology_metric_values(metric_id);
