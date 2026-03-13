package domain

import "time"

type AdminUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUserCreate struct {
	Username string
	Password string
}

