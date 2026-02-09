CREATE TABLE IF NOT EXISTS technology_sdgs (
    technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    sdg_id        UUID NOT NULL REFERENCES sdgs(id) ON DELETE CASCADE,
    PRIMARY KEY (technology_id, sdg_id)
    );
CREATE INDEX IF NOT EXISTS idx_technology_sdgs_sdg_id ON technology_sdgs(sdg_id);

CREATE TABLE IF NOT EXISTS technology_tags (
    technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    tag_id        UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (technology_id, tag_id)
    );
CREATE INDEX IF NOT EXISTS idx_technology_tags_tag_id ON technology_tags(tag_id);

CREATE TABLE IF NOT EXISTS technology_organizations (
    technology_id   UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    PRIMARY KEY (technology_id, organization_id)
    );
CREATE INDEX IF NOT EXISTS idx_technology_organizations_org_id ON technology_organizations(organization_id);



