CREATE TABLE IF NOT EXISTS technologies (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    trend_id UUID NOT NULL,

    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,

    readiness_level INTEGER NOT NULL,
    radar_index INTEGER NOT NULL,

    order_index INTEGER NOT NULL DEFAULT 0,

    custom_metric_1 NUMERIC,
    custom_metric_2 NUMERIC,
    custom_metric_3 NUMERIC,
    custom_metric_4 NUMERIC,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_technologies_trend
    FOREIGN KEY (trend_id)
    REFERENCES trends(id)
    ON DELETE RESTRICT,

    CONSTRAINT chk_readiness_level
    CHECK (readiness_level BETWEEN 1 AND 9),

    CONSTRAINT chk_radar_index
    CHECK (radar_index BETWEEN 1 AND 99)
    );

CREATE INDEX IF NOT EXISTS idx_technologies_trend_id
    ON technologies(trend_id);

CREATE INDEX IF NOT EXISTS idx_technologies_order_index
    ON technologies(order_index);
