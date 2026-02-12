CREATE TABLE IF NOT EXISTS trends (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 slug TEXT NOT NULL UNIQUE,
 name TEXT NOT NULL,
 description TEXT,
 image_url TEXT,
 order_index INTEGER NOT NULL DEFAULT 0 CHECK (order_index >= 0),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_trends_order_index ON trends(order_index);
CREATE INDEX IF NOT EXISTS idx_trends_slug ON trends(slug);
CREATE INDEX IF NOT EXISTS idx_trends_name_trgm ON trends USING gin (name gin_trgm_ops);

DROP TRIGGER IF EXISTS trg_trends_updated_at ON trends;
CREATE TRIGGER trg_trends_updated_at
    BEFORE UPDATE ON trends
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
