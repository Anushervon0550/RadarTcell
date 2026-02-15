BEGIN;

TRUNCATE
    technology_organizations,
    technology_tags,
    technology_sdgs,
    technologies,
    organizations,
    tags,
    sdgs,
    trends
    CASCADE;

INSERT INTO trends (slug, name, description, order_index)
VALUES
    ('ai', 'Artificial Intelligence', 'AI-related trends', 1),
    ('networks', 'Next-Gen Networks', 'Future networks trends', 2);

INSERT INTO sdgs (code, title)
VALUES
    ('SDG 09', 'Industry, Innovation and Infrastructure'),
    ('SDG 03', 'Good Health and Well-Being');

INSERT INTO tags (slug, title, category)
VALUES
    ('artificial-intelligence', 'Artificial Intelligence', 'Domain'),
    ('ml', 'Machine Learning', 'Domain'),
    ('telecom', 'Telecom', 'Industry');

INSERT INTO organizations (slug, name, logo_url, headquarters)
VALUES
    ('openai', 'OpenAI', 'https://example.com/openai.png', 'USA'),
    ('tcell', 'Tcell', 'https://example.com/tcell.png', 'Tajikistan');

-- Технологии (4 штуки для старта)
WITH t AS (
    SELECT id, slug FROM trends
),
     ins AS (
         INSERT INTO technologies (
                                   slug, list_index, name, description_short, readiness_level,
                                   custom_metric_1, custom_metric_2, custom_metric_3, custom_metric_4,
                                   trend_id
             )
             VALUES
                 ('edge-llm', 1, 'Edge LLM', 'LLM on device', 6, 0.7, 0.4, 0.8, 0.2, (SELECT id FROM t WHERE slug='ai')),
                 ('fraud-ml', 2, 'Fraud Detection ML', 'ML for fraud', 7, 0.6, 0.5, 0.4, 0.7, (SELECT id FROM t WHERE slug='ai')),
                 ('open-ran', 3, 'Open RAN', 'Open radio access network', 5, 0.5, 0.6, 0.3, 0.4, (SELECT id FROM t WHERE slug='networks')),
                 ('network-slicing', 4, 'Network Slicing', '5G slicing', 6, 0.4, 0.7, 0.6, 0.5, (SELECT id FROM t WHERE slug='networks'))
             RETURNING id, slug
     )
-- связи
INSERT INTO technology_sdgs (technology_id, sdg_id)
SELECT tech.id, s.id
FROM technologies tech
         JOIN sdgs s ON s.code IN ('SDG 09', 'SDG 03')
WHERE tech.slug IN ('edge-llm', 'fraud-ml');

INSERT INTO technology_tags (technology_id, tag_id)
SELECT tech.id, tag.id
FROM technologies tech
         JOIN tags tag ON (tech.slug='edge-llm' AND tag.slug IN ('artificial-intelligence','ml'))
    OR (tech.slug='open-ran' AND tag.slug IN ('telecom'));

INSERT INTO technology_organizations (technology_id, organization_id)
SELECT tech.id, org.id
FROM technologies tech
         JOIN organizations org ON (tech.slug IN ('edge-llm','fraud-ml') AND org.slug='openai')
    OR (tech.slug IN ('open-ran','network-slicing') AND org.slug='tcell');

COMMIT;
