package domain

import "strings"

type Technology struct {
	ID               string
	Slug             string
	Index            int
	Name             string
	DescriptionShort *string
	DescriptionFull  *string
	TRL              int

	CustomMetric1 *float64
	CustomMetric2 *float64
	CustomMetric3 *float64
	CustomMetric4 *float64

	ImageURL   *string
	SourceLink *string

	TrendID   string
	TrendSlug string
	TrendName string
}

type MetricRange struct {
	Min float64
	Max float64
}

type TechnologyListParams struct {
	Search         string
	TrendID        string
	SDGID          string
	TagID          string
	OrganizationID string

	TRLMin    int
	TRLMax    int
	HasTRLMin bool
	HasTRLMax bool

	SortBy string
	Order  string

	Page  int
	Limit int
	Cursor string

	Highlight []string
	OnlyIDs   []string // нужно для highlight-фильтрации
	Locale    string
}

type TechnologyViewBase struct {
	ID               string  `json:"id"`
	Slug             string  `json:"slug"`
	Index            int     `json:"index"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
	TRL              int     `json:"trl"`

	TrendID   string `json:"trend_id"`
	TrendSlug string `json:"trend_slug"`
	TrendName string `json:"trend_name"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty"`

	Angle  float64 `json:"angle"`
	Radius float64 `json:"radius"`
}

type TechnologyMetricValue struct {
	MetricID string   `json:"metric_id"`
	FieldKey *string  `json:"field_key,omitempty"`
	Value    *float64 `json:"value,omitempty"`
}

func MetricValueByFieldKey(items []TechnologyMetricValue, fieldKey string) *float64 {
	for _, it := range items {
		if it.FieldKey == nil || it.Value == nil {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(*it.FieldKey), fieldKey) {
			return it.Value
		}
	}
	return nil
}

type TechnologyListItem struct {
	TechnologyViewBase
	CustomMetrics []TechnologyMetricValue `json:"custom_metrics,omitempty"`

	CustomMetric1Norm float64 `json:"custom_metric_1_norm"`
	CustomMetric2Norm float64 `json:"custom_metric_2_norm"`
	CustomMetric3Norm float64 `json:"custom_metric_3_norm"`
	CustomMetric4Norm float64 `json:"custom_metric_4_norm"`
}

type TechnologyListResult struct {
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Total int                  `json:"total"`
	NextCursor string          `json:"next_cursor,omitempty"`
	Items []TechnologyListItem `json:"items"`
}

type TechnologyCard struct {
	TechnologyViewBase
	CustomMetrics []TechnologyMetricValue `json:"custom_metrics,omitempty"`

	DescriptionFull  *string `json:"description_full,omitempty"`

	ImageURL   *string `json:"image_url,omitempty"`
	SourceLink *string `json:"source_link,omitempty"`


	Tags          []Tag          `json:"tags"`
	SDGs          []SDG          `json:"sdgs"`
	Organizations []Organization `json:"organizations"`
}

type TechnologyCardData struct {
	Technology    Technology
	Tags          []Tag
	SDGs          []SDG
	Organizations []Organization
}

