package model

import "time"

type User struct {
	ID          int       `db:"id"`
	Username    string    `db:"username"`
	Description string    `db:"description"`
	Email       string    `db:"email"`
	Picture     string    `db:"picture"`
	Role        Role      `db:"role"`
	CreatedAt   time.Time `db:"created_at"`
}

type Role string

const (
	UserRole      Role = "user"
	ModeratorRole Role = "moderator"
	AdminRole     Role = "admin"
)
