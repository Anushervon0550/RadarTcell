package domain

// Роли администраторов. Сейчас система одноролевая (admin),
// но роль выносится в claim токена и проверяется middleware,
// чтобы можно было безопасно добавить новые роли позже.
const (
	RoleAdmin = "admin"
)

// Principal — аутентифицированный субъект запроса.
type Principal struct {
	Subject string
	Role    string
}
