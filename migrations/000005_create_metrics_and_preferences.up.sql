CREATE TABLE IF NOT EXISTS metrics_definitions (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name TEXT NOT NULL UNIQUE,
description TEXT,
type TEXT NOT NULL CHECK (type IN ('distance', 'bubble', 'bar')),
orderable BOOLEAN NOT NULL DEFAULT TRUE,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- user preferences: API ожидает settings-объект :contentReference[oaicite:7]{index=7}
CREATE TABLE IF NOT EXISTS user_preferences (
user_id TEXT PRIMARY KEY,
settings JSONB NOT NULL DEFAULT '{}'::jsonb,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

DROP TRIGGER IF EXISTS trg_metrics_updated_at ON metrics_definitions;
CREATE TRIGGER trg_metrics_updated_at
    BEFORE UPDATE ON metrics_definitions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_user_prefs_updated_at ON user_preferences;
CREATE TRIGGER trg_user_prefs_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- дефолтные метрики по смыслу ТЗ :contentReference[oaicite:8]{index=8}
INSERT INTO metrics_definitions (name, description, type, orderable)
VALUES
    ('Technology Readiness Level', 'NASA TRL 1–9 readiness level', 'distance', TRUE),
    ('Custom Metric 01', 'Example of a custom metric', 'bubble', TRUE),
    ('Custom Metric 02', 'Example of a custom metric', 'bubble', TRUE),
    ('Custom Metric 03', 'Example of a custom metric', 'bar', TRUE),
    ('Custom Metric 04', 'Example of a custom metric', 'bar', TRUE)
ON CONFLICT (name) DO NOTHING;
