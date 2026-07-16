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

-- Trends
INSERT INTO trends (slug, name, description, order_index)
VALUES
    ('ai', 'Artificial Intelligence', 'AI-related trends', 1),
    ('networks', 'Next-Gen Networks', 'Future networks trends', 2);

-- SDGs
INSERT INTO sdgs (code, title)
VALUES
    ('SDG 09', 'Industry, Innovation and Infrastructure'),
    ('SDG 03', 'Good Health and Well-Being');

-- Tags
INSERT INTO tags (slug, title, category)
VALUES
    ('artificial-intelligence', 'Artificial Intelligence', 'Domain'),
    ('ml', 'Machine Learning', 'Domain'),
    ('telecom', 'Telecom', 'Industry');

-- Organizations
INSERT INTO organizations (slug, name, logo_url, headquarters)
VALUES
    ('openai', 'OpenAI', 'https://example.com/openai.png', 'USA'),
    ('tcell', 'Tcell', 'https://example.com/tcell.png', 'Tajikistan');

-- Technologies (4 штуки для старта)
INSERT INTO technologies (
    slug,
    list_index,
    name,
    description_short,
    readiness_level,
    custom_metric_1,
    custom_metric_2,
    custom_metric_3,
    custom_metric_4,
    trend_id
)
VALUES
    (
        'edge-llm',
        1,
        'Edge LLM',
        'LLM on device',
        6,
        0.7, 0.4, 0.8, 0.2,
        (SELECT id FROM trends WHERE slug = 'ai')
    ),
    (
        'fraud-ml',
        2,
        'Fraud Detection ML',
        'ML for fraud',
        7,
        0.6, 0.5, 0.4, 0.7,
        (SELECT id FROM trends WHERE slug = 'ai')
    ),
    (
        'open-ran',
        3,
        'Open RAN',
        'Open radio access network',
        5,
        0.5, 0.6, 0.3, 0.4,
        (SELECT id FROM trends WHERE slug = 'networks')
    ),
    (
        'network-slicing',
        4,
        'Network Slicing',
        '5G slicing',
        6,
        0.4, 0.7, 0.6, 0.5,
        (SELECT id FROM trends WHERE slug = 'networks')
    );

-- Связи technologies <-> sdgs (точечные, без "перемножения")
INSERT INTO technology_sdgs (technology_id, sdg_id)
SELECT t.id, s.id
FROM technologies t
         JOIN sdgs s ON
    (t.slug = 'edge-llm'        AND s.code = 'SDG 09')
        OR (t.slug = 'fraud-ml'        AND s.code = 'SDG 03')
        OR (t.slug = 'open-ran'        AND s.code = 'SDG 09')
        OR (t.slug = 'network-slicing' AND s.code = 'SDG 09')
ON CONFLICT DO NOTHING;

-- Связи technologies <-> tags
INSERT INTO technology_tags (technology_id, tag_id)
SELECT t.id, tag.id
FROM technologies t
         JOIN tags tag ON
    (t.slug = 'edge-llm'        AND tag.slug IN ('artificial-intelligence', 'ml'))
        OR (t.slug = 'fraud-ml'        AND tag.slug IN ('ml'))
        OR (t.slug = 'open-ran'        AND tag.slug IN ('telecom'))
        OR (t.slug = 'network-slicing' AND tag.slug IN ('telecom'))
ON CONFLICT DO NOTHING;

-- Связи technologies <-> organizations
INSERT INTO technology_organizations (technology_id, organization_id)
SELECT t.id, org.id
FROM technologies t
         JOIN organizations org ON
    (t.slug IN ('edge-llm', 'fraud-ml') AND org.slug = 'openai')
        OR (t.slug IN ('open-ran', 'network-slicing') AND org.slug = 'tcell')
ON CONFLICT DO NOTHING;

COMMIT;