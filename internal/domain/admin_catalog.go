package domain

type TrendUpsert struct {
	Slug        string // required for create
	Name        string
	Description *string
	ImageURL    *string
	Order       int
}

// AdminTrend представляет полный объект тренда для админ‑GET.
type AdminTrend struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	OrderIndex  int     `json:"order_index"`
}

// AdminSDG представляет полный объект SDG для админ‑GET.
type AdminSDG struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type TagUpsert struct {
	Slug        string // required for create
	Title       string
	Category    string
	Description *string
}
