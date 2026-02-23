package domain

type OrganizationUpsert struct {
	Slug    string // required for create
	Name    string
	LogoURL *string
}

type MetricDefinitionUpsert struct {
	Name        string
	Type        string
	Description *string
	Orderable   bool
	FieldKey    *string
}
