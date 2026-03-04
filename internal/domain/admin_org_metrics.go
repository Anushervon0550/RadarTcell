package domain

type OrganizationUpsert struct {
	Slug         string // required for create
	Name         string
	LogoURL      *string
	Description  *string
	Website      *string
	Headquarters *string
}

type MetricDefinitionUpsert struct {
	Name        string
	Type        string
	Description *string
	Orderable   bool
	FieldKey    *string
}
