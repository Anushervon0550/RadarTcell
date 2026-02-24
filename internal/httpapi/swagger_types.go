package httpapi

type ErrorResponse struct {
	Error string `json:"error" example:"not found"`
}

type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

type IDResponse struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type IDSlugResponse struct {
	ID   string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug string `json:"slug" example:"edge-llm"`
}

type AdminLoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
}

type AdminLoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AdminMeResponse struct {
	Role string `json:"role" example:"admin"`
	User string `json:"user" example:"admin"`
}

type TrendUpsertRequest struct {
	Slug        string  `json:"slug,omitempty" example:"ai"`
	Name        string  `json:"name" example:"Artificial Intelligence"`
	Description *string `json:"description,omitempty" example:"AI-related trends"`
	OrderIndex  int     `json:"order_index" example:"1"`
}

type TagUpsertRequest struct {
	Slug        string  `json:"slug,omitempty" example:"ml"`
	Title       string  `json:"title" example:"Machine Learning"`
	Category    string  `json:"category" example:"Domain"`
	Description *string `json:"description,omitempty" example:"ML technologies"`
}

type OrganizationUpsertRequest struct {
	Slug         string  `json:"slug,omitempty" example:"openai"`
	Name         string  `json:"name" example:"OpenAI"`
	LogoURL      *string `json:"logo_url,omitempty" example:"https://example.com/openai.png"`
	Headquarters *string `json:"headquarters,omitempty" example:"USA"`
	Description  *string `json:"description,omitempty" example:"AI research company"`
}

type MetricUpsertRequest struct {
	Name        string  `json:"name" example:"Custom Metric 01"`
	Type        string  `json:"type" example:"bar"` // bubble|bar|distance
	Description *string `json:"description,omitempty" example:"Example metric"`
	Orderable   bool    `json:"orderable" example:"true"`
	FieldKey    *string `json:"field_key,omitempty" example:"custom_metric_1"`
}

type TechnologyUpsertRequest struct {
	Slug string `json:"slug,omitempty" example:"edge-llm"`

	Index int    `json:"index" example:"1"`
	Name  string `json:"name" example:"Edge LLM"`
	TRL   int    `json:"trl" example:"6"`

	TrendSlug string `json:"trend_slug" example:"ai"`

	DescriptionShort *string `json:"description_short,omitempty" example:"LLM on device"`
	DescriptionFull  *string `json:"description_full,omitempty" example:"Detailed description..."`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty" example:"0.7"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty" example:"0.4"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty" example:"0.8"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty" example:"0.2"`

	ImageURL   *string `json:"image_url,omitempty" example:"https://example.com/img.png"`
	SourceLink *string `json:"source_link,omitempty" example:"https://example.com/article"`

	TagSlugs          []string `json:"tag_slugs,omitempty" example:"ml,artificial-intelligence"`
	SDGCodes          []string `json:"sdg_codes,omitempty" example:"SDG 09,SDG 03"`
	OrganizationSlugs []string `json:"organization_slugs,omitempty" example:"openai,tcell"`
}
type TrendDTO struct {
	ID                string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug              string `json:"slug" example:"ai"`
	Name              string `json:"name" example:"Artificial Intelligence"`
	TechnologiesCount int    `json:"technologies_count" example:"2"`
}

type SDGDTO struct {
	ID                string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Code              string `json:"code" example:"SDG 09"`
	Title             string `json:"title" example:"Industry, Innovation and Infrastructure"`
	TechnologiesCount int    `json:"technologies_count" example:"3"`
}

type TagDTO struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug     string `json:"slug" example:"ml"`
	Title    string `json:"title" example:"Machine Learning"`
	Category string `json:"category" example:"Domain"`
}

type OrganizationDTO struct {
	ID                string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug              string  `json:"slug" example:"openai"`
	Name              string  `json:"name" example:"OpenAI"`
	LogoURL           *string `json:"logo_url,omitempty" example:"https://example.com/openai.png"`
	TechnologiesCount int     `json:"technologies_count" example:"2"`
}

type MetricDTO struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string  `json:"name" example:"Custom Metric 01"`
	Type        string  `json:"type" example:"bubble"`
	Description *string `json:"description,omitempty" example:"Example of a custom metric"`
	Orderable   bool    `json:"orderable" example:"true"`
	FieldKey    *string `json:"field_key,omitempty" example:"custom_metric_1"`
}

type TechnologyListItemDTO struct {
	ID               string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug             string  `json:"slug" example:"edge-llm"`
	Index            int     `json:"index" example:"1"`
	Name             string  `json:"name" example:"Edge LLM"`
	DescriptionShort *string `json:"description_short,omitempty" example:"LLM on device"`
	TRL              int     `json:"trl" example:"6"`

	TrendID   string `json:"trend_id" example:"550e8400-e29b-41d4-a716-446655440111"`
	TrendSlug string `json:"trend_slug" example:"ai"`
	TrendName string `json:"trend_name" example:"Artificial Intelligence"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty" example:"0.7"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty" example:"0.4"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty" example:"0.8"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty" example:"0.2"`
}

type TechnologyDetailDTO struct {
	ID               string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug             string  `json:"slug" example:"edge-llm"`
	Index            int     `json:"index" example:"1"`
	Name             string  `json:"name" example:"Edge LLM"`
	DescriptionShort *string `json:"description_short,omitempty" example:"LLM on device"`
	DescriptionFull  *string `json:"description_full,omitempty" example:"Detailed description"`
	TRL              int     `json:"trl" example:"6"`

	TrendID   string `json:"trend_id" example:"550e8400-e29b-41d4-a716-446655440111"`
	TrendSlug string `json:"trend_slug" example:"ai"`
	TrendName string `json:"trend_name" example:"Artificial Intelligence"`

	CustomMetric1 *float64 `json:"custom_metric_1,omitempty" example:"0.7"`
	CustomMetric2 *float64 `json:"custom_metric_2,omitempty" example:"0.4"`
	CustomMetric3 *float64 `json:"custom_metric_3,omitempty" example:"0.8"`
	CustomMetric4 *float64 `json:"custom_metric_4,omitempty" example:"0.2"`

	ImageURL   *string `json:"image_url,omitempty" example:"https://example.com/img.png"`
	SourceLink *string `json:"source_link,omitempty" example:"https://example.com/article"`

	Tags          []TagDTO          `json:"tags,omitempty"`
	SDGs          []SDGDTO          `json:"sdgs,omitempty"`
	Organizations []OrganizationDTO `json:"organizations,omitempty"`
}

type TechnologyListResponse struct {
	Page  int                     `json:"page" example:"1"`
	Limit int                     `json:"limit" example:"20"`
	Total int                     `json:"total" example:"2"`
	Items []TechnologyListItemDTO `json:"items"`
}

type MetricValueResponse struct {
	FieldKey     string  `json:"field_key" example:"custom_metric_1"`
	MetricID     string  `json:"metric_id" example:"26dea5c3-6eb3-407a-927c-414c787c0cdf"`
	MetricName   string  `json:"metric_name" example:"Custom Metric 01"`
	TechnologyID string  `json:"technology_id" example:"5ab5c441-465e-423b-a940-b47fbfaf088b"`
	Type         string  `json:"type" example:"bubble"`
	Value        float64 `json:"value" example:"0.7"`
}

type PreferencesSaveRequest struct {
	UserID   string         `json:"user_id" example:"u1"`
	Settings map[string]any `json:"settings"`
}

type PreferencesSaveResponse struct {
	Status string `json:"status" example:"ok"`
}

type PreferencesGetResponse struct {
	UserID   string         `json:"user_id" example:"u1"`
	Settings map[string]any `json:"settings"`
}
