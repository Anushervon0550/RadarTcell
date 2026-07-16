DROP TABLE IF EXISTS technology_organizations;
DROP TABLE IF EXISTS technology_tags;
DROP TABLE IF EXISTS technology_sdgs;

DROP TRIGGER IF EXISTS trg_orgs_updated_at ON organizations;
DROP TRIGGER IF EXISTS trg_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS trg_sdgs_updated_at ON sdgs;

DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS sdgs;
