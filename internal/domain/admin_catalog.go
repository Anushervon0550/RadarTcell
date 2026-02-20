package domain

type TrendUpsert struct {
	Slug  string // required for create
	Name  string
	Order int
}

type TagUpsert struct {
	Slug        string // required for create
	Title       string
	Category    string
	Description *string
}
