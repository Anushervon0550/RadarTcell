CREATE TABLE IF NOT EXISTS technologies (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
slug TEXT NOT NULL UNIQUE,
list_index SMALLINT NOT NULL CHECK (list_index BETWEEN 1 AND 99),
name TEXT NOT NULL,
description_short TEXT,
description_full TEXT,
readiness_level SMALLINT NOT NULL CHECK (readiness_level BETWEEN 1 AND 9),
custom_metric_1 NUMERIC,
custom_metric_2 NUMERIC,
custom_metric_3 NUMERIC,
custom_metric_4 NUMERIC,
image_url TEXT,
source_link TEXT,
trend_id UUID NOT NULL REFERENCES trends(id) ON DELETE RESTRICT,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- индексы из ТЗ для производительности /api/technologies :contentReference[oaicite:4]{index=4}
CREATE INDEX IF NOT EXISTS idx_technologies_trend_id ON technologies(trend_id);
CREATE INDEX IF NOT EXISTS idx_technologies_readiness_level ON technologies(readiness_level);
CREATE INDEX IF NOT EXISTS idx_technologies_custom_metric_1 ON technologies(custom_metric_1);
CREATE INDEX IF NOT EXISTS idx_technologies_custom_metric_2 ON technologies(custom_metric_2);
CREATE INDEX IF NOT EXISTS idx_technologies_custom_metric_3 ON technologies(custom_metric_3);
CREATE INDEX IF NOT EXISTS idx_technologies_custom_metric_4 ON technologies(custom_metric_4);

-- поиск по имени
CREATE INDEX IF NOT EXISTS idx_technologies_name_trgm ON technologies USING gin (name gin_trgm_ops);

DROP TRIGGER IF EXISTS trg_technologies_updated_at ON technologies;
CREATE TRIGGER trg_technologies_updated_at
    BEFORE UPDATE ON technologies
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
