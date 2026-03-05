package domain

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

	Highlight []string
	OnlyIDs   []string // нужно для highlight-фильтрации
	Locale    string
}

type TechnologyListItem struct {
	ID               string  `json:"id"`
	Slug             string  `json:"slug"`
	Index            int     `json:"index"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
	TRL              int     `json:"trl"`
	Angle            float64 `json:"angle"`
	Radius           float64 `json:"radius"`

	TrendID   string `json:"trend_id"`
	TrendSlug string `json:"trend_slug"`
	TrendName string `json:"trend_name"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty"`

	CustomMetric1Norm float64 `json:"custom_metric_1_norm"`
	CustomMetric2Norm float64 `json:"custom_metric_2_norm"`
	CustomMetric3Norm float64 `json:"custom_metric_3_norm"`
	CustomMetric4Norm float64 `json:"custom_metric_4_norm"`
}

type TechnologyListResult struct {
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
	Total int                  `json:"total"`
	Items []TechnologyListItem `json:"items"`
}
type TechnologyCard struct {
	ID               string  `json:"id"`
	Slug             string  `json:"slug"`
	Index            int     `json:"index"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
	DescriptionFull  *string `json:"description_full,omitempty"`
	TRL              int     `json:"trl"`

	TrendID   string `json:"trend_id"`
	TrendSlug string `json:"trend_slug"`
	TrendName string `json:"trend_name"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty"`

	ImageURL   *string `json:"image_url,omitempty"`
	SourceLink *string `json:"source_link,omitempty"`

	Angle  float64 `json:"angle"`
	Radius float64 `json:"radius"`

	Tags          []Tag          `json:"tags"`
	SDGs          []SDG          `json:"sdgs"`
	Organizations []Organization `json:"organizations"`
}
