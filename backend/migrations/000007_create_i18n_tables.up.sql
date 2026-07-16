CREATE TABLE IF NOT EXISTS trend_i18n (
    trend_id UUID NOT NULL REFERENCES trends(id) ON DELETE CASCADE,
    locale TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    PRIMARY KEY (trend_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_trend_i18n_locale ON trend_i18n(locale);

CREATE TABLE IF NOT EXISTS technology_i18n (
    technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    locale TEXT NOT NULL,
    name TEXT NOT NULL,
    description_short TEXT,
    description_full TEXT,
    PRIMARY KEY (technology_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_technology_i18n_locale ON technology_i18n(locale);

CREATE TABLE IF NOT EXISTS metric_definition_i18n (
    metric_id UUID NOT NULL REFERENCES metrics_definitions(id) ON DELETE CASCADE,
    locale TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    PRIMARY KEY (metric_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_metric_definition_i18n_locale ON metric_definition_i18n(locale);

