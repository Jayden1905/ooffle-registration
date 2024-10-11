package types

import (
	"context"
	"time"

	"github.com/jayden1905/event-registration-software/cmd/pkg/database"
)

type User struct {
	ID        int32     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore interface {
	GetUserByEmail(email string) (*database.User, error)
	GetUserByID(id int32) (*database.User, error)
	CreateUser(ctx context.Context, user *database.User) error
	CreateSuperUser(ctx context.Context, user *database.User) error
	UpdateUserToSuperUser(ctx context.Context, id int32) error
	UpdateUserToNormalUser(ctx context.Context, id int32) error
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=20"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func DatabaseUserToUser(u *database.User) *User {
	return &User{
		ID:        int32(u.UserID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
