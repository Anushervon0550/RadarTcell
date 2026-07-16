package domain

import "time"

type TechnologyMetricValueUpsert struct {
	MetricID string
	Value    *float64
}

type TechnologyUpsert struct {
	Slug string // required for create

	Index int
	Name  string
	TRL   int

	TrendSlug string

	DescriptionShort *string
	DescriptionFull  *string

	CustomMetric1 *float64
	CustomMetric2 *float64
	CustomMetric3 *float64
	CustomMetric4 *float64

	ImageURL   *string
	SourceLink *string

	TagSlugs          []string
	SDGCodes          []string
	OrganizationSlugs []string
	CustomMetrics     []TechnologyMetricValueUpsert
}

type TechnologyAdmin struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`

	Index int    `json:"index"`
	Name  string `json:"name"`
	TRL   int    `json:"trl"`

	TrendSlug string `json:"trend_slug"`
	TrendName string `json:"trend_name,omitempty"`

	DescriptionShort *string `json:"description_short,omitempty"`
	DescriptionFull  *string `json:"description_full,omitempty"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty"`

	ImageURL   *string `json:"image_url,omitempty"`
	SourceLink *string `json:"source_link,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	CustomMetrics []TechnologyMetricValue `json:"custom_metrics,omitempty"`

	TagSlugs          []string `json:"tag_slugs,omitempty"`
	SDGCodes          []string `json:"sdg_codes,omitempty"`
	OrganizationSlugs []string `json:"organization_slugs,omitempty"`
}

type AdminTechnologyListParams struct {
	Page  int
	Limit int
	IncludeDeleted bool
}

type AdminTechnologyListResult struct {
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Total int                `json:"total"`
	Items []TechnologyAdmin  `json:"items"`
}

