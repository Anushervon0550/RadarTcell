BEGIN;

-- очищаем только технологии и связи (каталог не трогаем)
TRUNCATE
    technology_organizations,
    technology_tags,
    technology_sdgs,
    technologies
    CASCADE;

WITH
    t_ai  AS (SELECT id FROM trends WHERE slug = 'ai'),
    t_net AS (SELECT id FROM trends WHERE slug = 'networks'),

    tag_ai  AS (SELECT id FROM tags WHERE slug = 'artificial-intelligence'),
    tag_ml  AS (SELECT id FROM tags WHERE slug = 'ml'),
    tag_tel AS (SELECT id FROM tags WHERE slug = 'telecom'),

    sdg09 AS (SELECT id FROM sdgs WHERE code = 'SDG 09'),
    sdg03 AS (SELECT id FROM sdgs WHERE code = 'SDG 03'),

    org_openai AS (SELECT id FROM organizations WHERE slug = 'openai'),
    org_tcell   AS (SELECT id FROM organizations WHERE slug = 'tcell'),

    ins AS (
        INSERT INTO technologies (
                                  slug, list_index, name, description_short, readiness_level,
                                  custom_metric_1, custom_metric_2, custom_metric_3, custom_metric_4,
                                  trend_id
            )
            VALUES
                -- AI (8)
                ('edge-llm',            1,  'Edge LLM',             'LLM on device',             6, 0.70, 0.40, 0.80, 0.20, (SELECT id FROM t_ai)),
                ('fraud-ml',            2,  'Fraud Detection ML',   'ML for fraud',              7, 0.60, 0.50, 0.40, 0.70, (SELECT id FROM t_ai)),
                ('rag-platform',        3,  'RAG Platform',         'Retrieval + generation',    6, 0.65, 0.55, 0.75, 0.35, (SELECT id FROM t_ai)),
                ('ai-ops',              4,  'AIOps',                'ML for ops',                5, 0.55, 0.60, 0.50, 0.45, (SELECT id FROM t_ai)),
                ('customer-churn-ml',   5,  'Customer Churn ML',    'Predict churn',             6, 0.58, 0.62, 0.52, 0.48, (SELECT id FROM t_ai)),
                ('voice-analytics',     6,  'Voice Analytics',      'Speech insights',           5, 0.50, 0.45, 0.60, 0.55, (SELECT id FROM t_ai)),
                ('recommendation-ai',   7,  'Recommendation AI',    'Personalization',           7, 0.72, 0.66, 0.70, 0.40, (SELECT id FROM t_ai)),
                ('anomaly-detection',   8,  'Anomaly Detection',    'Detect anomalies',          6, 0.62, 0.58, 0.68, 0.42, (SELECT id FROM t_ai)),

                -- Networks (8)
                ('open-ran',            9,  'Open RAN',             'Open radio access network', 5, 0.50, 0.60, 0.30, 0.40, (SELECT id FROM t_net)),
                ('network-slicing',     10, 'Network Slicing',      '5G slicing',                6, 0.40, 0.70, 0.60, 0.50, (SELECT id FROM t_net)),
                ('private-5g',          11, 'Private 5G',            'Enterprise 5G',             6, 0.48, 0.64, 0.55, 0.52, (SELECT id FROM t_net)),
                ('edge-mec',            12, 'MEC Edge',              'Edge compute for telco',    5, 0.46, 0.57, 0.53, 0.49, (SELECT id FROM t_net)),
                ('wifi-7',              13, 'Wi-Fi 7',               'Next gen Wi-Fi',            7, 0.60, 0.50, 0.65, 0.45, (SELECT id FROM t_net)),
                ('6g-research',         14, '6G Research',           'Future 6G direction',       3, 0.30, 0.35, 0.25, 0.20, (SELECT id FROM t_net)),
                ('network-digital-twin',15, 'Network Digital Twin',  'Simulate networks',         4, 0.38, 0.42, 0.40, 0.33, (SELECT id FROM t_net)),
                ('son-automation',      16, 'SON Automation',        'Self-organizing network',   6, 0.52, 0.55, 0.58, 0.47, (SELECT id FROM t_net))
            RETURNING id, slug
    ),

-- SDG links
    ins_sdg_09 AS (
        INSERT INTO technology_sdgs (technology_id, sdg_id)
            SELECT tech.id, (SELECT id FROM sdg09)
            FROM ins tech
            WHERE tech.slug IN (
                                'edge-llm','fraud-ml','rag-platform','ai-ops','customer-churn-ml','anomaly-detection',
                                'open-ran','network-slicing','private-5g','edge-mec','wifi-7','network-digital-twin','son-automation'
                )
            RETURNING 1
    ),
    ins_sdg_03 AS (
        INSERT INTO technology_sdgs (technology_id, sdg_id)
            SELECT tech.id, (SELECT id FROM sdg03)
            FROM ins tech
            WHERE tech.slug IN ('voice-analytics','recommendation-ai')
            RETURNING 1
    ),

-- TAG links
    ins_tag_ai AS (
        INSERT INTO technology_tags (technology_id, tag_id)
            SELECT tech.id, (SELECT id FROM tag_ai)
            FROM ins tech
            WHERE tech.slug IN ('edge-llm','rag-platform','recommendation-ai')
            RETURNING 1
    ),
    ins_tag_ml AS (
        INSERT INTO technology_tags (technology_id, tag_id)
            SELECT tech.id, (SELECT id FROM tag_ml)
            FROM ins tech
            WHERE tech.slug IN ('edge-llm','fraud-ml','rag-platform','ai-ops','customer-churn-ml','voice-analytics','recommendation-ai','anomaly-detection')
            RETURNING 1
    ),
    ins_tag_tel AS (
        INSERT INTO technology_tags (technology_id, tag_id)
            SELECT tech.id, (SELECT id FROM tag_tel)
            FROM ins tech
            WHERE tech.slug IN ('open-ran','network-slicing','private-5g','edge-mec','wifi-7','6g-research','network-digital-twin','son-automation')
            RETURNING 1
    ),

-- ORG links
    ins_org_openai AS (
        INSERT INTO technology_organizations (technology_id, organization_id)
            SELECT tech.id, (SELECT id FROM org_openai)
            FROM ins tech
            WHERE tech.slug IN ('edge-llm','fraud-ml','rag-platform','ai-ops','customer-churn-ml','voice-analytics','recommendation-ai','anomaly-detection')
            RETURNING 1
    ),
    ins_org_tcell AS (
        INSERT INTO technology_organizations (technology_id, organization_id)
            SELECT tech.id, (SELECT id FROM org_tcell)
            FROM ins tech
            WHERE tech.slug IN ('open-ran','network-slicing','private-5g','edge-mec','wifi-7','6g-research','network-digital-twin','son-automation')
            RETURNING 1
    )

SELECT 1;

COMMIT;