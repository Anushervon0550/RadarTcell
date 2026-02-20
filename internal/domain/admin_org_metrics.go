package domain

type OrganizationUpsert struct {
	Slug    string // required for create
	Name    string
	LogoURL *string
}

type MetricDefinitionUpsert struct {
	Name        string
	Type        string // "bubble" | "bar"
	Description *string
	Orderable   bool
}
