CREATE TABLE IF NOT EXISTS technologies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,

    -- по ТЗ
    "index" INTEGER NOT NULL CHECK ("index" BETWEEN 1 AND 99),
    description_short TEXT,
    description_full  TEXT,
    readiness_level INTEGER NOT NULL CHECK (readiness_level BETWEEN 1 AND 9),

    custom_metric_1 NUMERIC,
    custom_metric_2 NUMERIC,
    custom_metric_3 NUMERIC,
    custom_metric_4 NUMERIC,

    image_url   TEXT,
    source_link TEXT,

    trend_id UUID NOT NULL REFERENCES trends(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

-- индексы для фильтрации/сортировки
CREATE INDEX IF NOT EXISTS idx_technologies_trend_id ON technologies(trend_id);
CREATE INDEX IF NOT EXISTS idx_technologies_index ON technologies("index");
CREATE INDEX IF NOT EXISTS idx_technologies_readiness_level ON technologies(readiness_level);

-- быстрый поиск по name (search=...)
CREATE INDEX IF NOT EXISTS idx_technologies_name_trgm
    ON technologies USING gin (name gin_trgm_ops);
