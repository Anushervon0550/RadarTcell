package domain

type Trend struct {
	ID                string `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	TechnologiesCount int    `json:"technologies_count"`
}

type SDG struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Title             string `json:"title"`
	TechnologiesCount int    `json:"technologies_count"`
}

type Tag struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Category    *string `json:"category,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Organization struct {
	ID                string  `json:"id"`
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	LogoURL           *string `json:"logo_url,omitempty"`
	Description       *string `json:"description,omitempty"`
	Website           *string `json:"website,omitempty"`
	Headquarters      *string `json:"headquarters,omitempty"`
	TechnologiesCount int     `json:"technologies_count"`
}

type MetricDefinition struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // distance | bubble | bar
	Description *string `json:"description,omitempty"`
	Orderable   bool    `json:"orderable"`
	FieldKey    *string `json:"field_key,omitempty"`
}
