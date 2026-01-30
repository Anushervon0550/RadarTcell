CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE trends(
    id UUID PRIMARY KEY DEFAULT gen_rendom_uuid(),
    slug TEXT NOT NULL UNIQUE ,
    name TEXT NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL DEFAULT 0,
    create_at TIMESTAMPZ NOT NULL DEFAULT now(),
    updated_up TIMESTAMPZ NOT NULL DEFAULT now()
):