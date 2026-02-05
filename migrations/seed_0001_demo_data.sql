-- DEV SEED (safe order: parents -> children)

BEGIN;

-- 0) чистим всё чтобы можно было запускать seed много раз
TRUNCATE
    technology_metric_values,
  technology_organizations,
  technology_sdgs,
  technology_tags,
  technologies,
  metrics_definitions,
  organizations,
  sdgs,
  tags,
  trends
CASCADE;

-- 1) TRENDS (родитель)
INSERT INTO trends (slug, name, description, image_url, order_index)
VALUES
    ('ai', 'AI', 'Artificial Intelligence', NULL, 1),
    ('energy', 'Energy', 'Energy & Sustainability', NULL, 2);

-- 2) SDGS
INSERT INTO sdgs (code, title, description, icon)
VALUES
    ('SDG09', 'Industry, Innovation and Infrastructure', NULL, NULL),
    ('SDG07', 'Affordable and Clean Energy', NULL, NULL),
    ('SDG13', 'Climate Action', NULL, NULL);

-- 3) TAGS
INSERT INTO tags (slug, title, category, description)
VALUES
    ('genai', 'GenAI', 'ai', NULL),
    ('storage', 'Energy Storage', 'energy', NULL),
    ('computer-vision', 'Computer Vision', 'ai', NULL);

-- 4) ORGS
INSERT INTO organizations (slug, name, website, headquarters)
VALUES
    ('openai', 'OpenAI', 'https://openai.com', 'San Francisco'),
    ('tesla', 'Tesla', 'https://www.tesla.com', 'Austin');

-- 5) TECHNOLOGIES (дети, требуют trends)
WITH t AS (
    SELECT
        (SELECT id FROM trends WHERE slug='ai')     AS ai_id,
        (SELECT id FROM trends WHERE slug='energy') AS en_id
)
INSERT INTO technologies
(trend_id, slug, name, "index", description_short, description_full,
 readiness_level, image_url, source_link,
 custom_metric_1, custom_metric_2, custom_metric_3, custom_metric_4)
VALUES
  ((SELECT ai_id FROM t), 'llm', 'Large Language Models', 10, 'LLM базовые модели', NULL, 7, NULL, NULL, 0.7, 0.5, NULL, NULL),
  ((SELECT ai_id FROM t), 'rag', 'RAG', 20, 'Retrieval Augmented Generation', NULL, 6, NULL, NULL, 0.6, 0.7, NULL, NULL),
  ((SELECT ai_id FROM t), 'cv', 'Computer Vision', 30, 'Vision models', NULL, 8, NULL, NULL, 0.8, 0.4, NULL, NULL),
  ((SELECT en_id FROM t), 'batteries', 'Advanced Batteries', 15, 'Новые аккумуляторы', NULL, 8, NULL, NULL, 0.7, 0.8, NULL, NULL),
  ((SELECT en_id FROM t), 'grid', 'Smart Grids', 25, 'Умные сети', NULL, 7, NULL, NULL, 0.5, 0.6, NULL, NULL),
  ((SELECT en_id FROM t), 'hydrogen', 'Green Hydrogen', 35, 'Зелёный водород', NULL, 5, NULL, NULL, 0.4, 0.7, NULL, NULL);

-- 6) LINKS: technology_tags
INSERT INTO technology_tags (technology_id, tag_id)
SELECT tech.id, tag.id
FROM technologies tech
         JOIN tags tag ON tag.slug IN ('genai','computer-vision')
WHERE tech.slug IN ('llm','rag','cv');

INSERT INTO technology_tags (technology_id, tag_id)
SELECT tech.id, tag.id
FROM technologies tech
         JOIN tags tag ON tag.slug IN ('storage')
WHERE tech.slug IN ('batteries','grid','hydrogen');

-- 7) LINKS: technology_organizations
INSERT INTO technology_organizations (technology_id, organization_id)
SELECT tech.id, org.id
FROM technologies tech
         JOIN organizations org ON org.slug='openai'
WHERE tech.slug IN ('llm','rag');

INSERT INTO technology_organizations (technology_id, organization_id)
SELECT tech.id, org.id
FROM technologies tech
         JOIN organizations org ON org.slug='tesla'
WHERE tech.slug IN ('batteries','grid');

-- 8) LINKS: technology_sdgs
INSERT INTO technology_sdgs (technology_id, sdg_id)
SELECT tech.id, s.id
FROM technologies tech
         JOIN sdgs s ON s.code IN ('SDG09','SDG13')
WHERE tech.slug IN ('llm','rag','cv');

INSERT INTO technology_sdgs (technology_id, sdg_id)
SELECT tech.id, s.id
FROM technologies tech
         JOIN sdgs s ON s.code IN ('SDG07','SDG13')
WHERE tech.slug IN ('batteries','grid','hydrogen');

-- 9) METRICS DEFINITIONS + VALUES
INSERT INTO metrics_definitions (name, description, type, orderable)
VALUES
    ('impact', 'Impact score', 'bar', true),
    ('cost', 'Cost score', 'distance', true);

INSERT INTO technology_metric_values (technology_id, metric_id, value)
SELECT tech.id, m.id, v.val
FROM (VALUES
          ('llm', 0.8),
          ('rag', 0.7),
          ('batteries', 0.9),
          ('hydrogen', 0.6)
     ) v(slug, val)
         JOIN technologies tech ON tech.slug = v.slug
         JOIN metrics_definitions m ON m.name = 'impact';

COMMIT;
