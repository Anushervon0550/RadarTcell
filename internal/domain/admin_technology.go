package domain

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
}
