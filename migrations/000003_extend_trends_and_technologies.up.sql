CREATE TABLE IF NOT EXISTS sdgs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    code TEXT NOT NULL UNIQUE,       -- например "SDG 09"
    title TEXT NOT NULL,
    description TEXT,
    icon TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_sdgs_code ON sdgs(code);
