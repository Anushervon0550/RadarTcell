CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS trends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    image_url TEXT,

    -- позиция тренда на круговой диаграмме
    order_index INTEGER NOT NULL DEFAULT 0 CHECK (order_index >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_trends_order_index ON trends(order_index);
CREATE INDEX IF NOT EXISTS idx_trends_slug ON trends(slug);

-- для поиска по имени (если будешь делать ILIKE/contains)
CREATE INDEX IF NOT EXISTS idx_trends_name_trgm
    ON trends USING gin (name gin_trgm_ops);
