package domain

type TrendI18nUpsert struct {
	Locale      string
	Name        string
	Description *string
}

type TechnologyI18nUpsert struct {
	Locale           string
	Name             string
	DescriptionShort *string
	DescriptionFull  *string
}

type MetricI18nUpsert struct {
	Locale      string
	Name        string
	Description *string
}

type TrendI18n struct {
	Locale      string  `json:"locale"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type TechnologyI18n struct {
	Locale           string  `json:"locale"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
	DescriptionFull  *string `json:"description_full,omitempty"`
}

type MetricI18n struct {
	Locale      string  `json:"locale"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
