CREATE TABLE IF NOT EXISTS sdgs (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
code TEXT NOT NULL UNIQUE,
title TEXT NOT NULL,
description TEXT,
icon TEXT,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tags (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
slug TEXT NOT NULL UNIQUE,
title TEXT NOT NULL,
category TEXT,
description TEXT,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS organizations (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
slug TEXT NOT NULL UNIQUE,
name TEXT NOT NULL,
logo_url TEXT,
description TEXT,
website TEXT,
headquarters TEXT,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- связи many-to-many по ТЗ :contentReference[oaicite:5]{index=5} :contentReference[oaicite:6]{index=6}
CREATE TABLE IF NOT EXISTS technology_sdgs (
technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
sdg_id UUID NOT NULL REFERENCES sdgs(id) ON DELETE CASCADE,
PRIMARY KEY (technology_id, sdg_id)
);

CREATE TABLE IF NOT EXISTS technology_tags (
technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
PRIMARY KEY (technology_id, tag_id)
);

CREATE TABLE IF NOT EXISTS technology_organizations (
technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
PRIMARY KEY (technology_id, organization_id)
);

-- индексы для обратных выборок (фильтры по sdg/tag/org)
CREATE INDEX IF NOT EXISTS idx_tech_sdgs_sdg_id ON technology_sdgs(sdg_id);
CREATE INDEX IF NOT EXISTS idx_tech_tags_tag_id ON technology_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_tech_orgs_org_id ON technology_organizations(organization_id);

-- триггеры updated_at
DROP TRIGGER IF EXISTS trg_sdgs_updated_at ON sdgs;
CREATE TRIGGER trg_sdgs_updated_at
    BEFORE UPDATE ON sdgs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_tags_updated_at ON tags;
CREATE TRIGGER trg_tags_updated_at
    BEFORE UPDATE ON tags
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_orgs_updated_at ON organizations;
CREATE TRIGGER trg_orgs_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
