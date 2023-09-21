package config

import (
	"time"
)

type DbUrl struct {
	DB_NAME string
	DB_PORT string
	DB_PASS string
	DB_USER string
	DB_HOST string
}

type User struct {
	Phone     uint64    `json:"phone"`
	Name      string    `json:"name"`
	Block     bool      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDelete struct {
	Phone uint64 `json:"phone"`
}
