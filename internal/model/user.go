package model

import (
	"time"

	"github.com/Newella-HQ/protos/gen/go/user"
)

type User struct {
	ID          string    `db:"id"`
	Username    string    `db:"username"`
	RealName    string    `db:"real_name"`
	Description string    `db:"description"`
	Email       string    `db:"email"`
	Picture     string    `db:"picture"`
	Role        Role      `db:"role"`
	CreatedAt   time.Time `db:"created_at"`
}

type Role string

func (r *Role) ToProto() user.Role {
	if r == nil {
		return -1
	}
	return user.Role(user.Role_value[string(*r)])
}

const (
	UserRole      Role = "user"
	ModeratorRole Role = "moderator"
	AdminRole     Role = "admin"
)
