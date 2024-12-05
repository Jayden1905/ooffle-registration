package types

import (
	"context"
	"time"
)

type User struct {
	ID           int32     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	Subscription string    `json:"subscription"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Verify       bool      `json:"verify"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStore interface {
	GetUsersPaginated(page int32, pageSize int32) ([]*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int32) (*User, error)
	GetUserRoleByID(id int32) (string, error)
	CreateUser(ctx context.Context, user *User) error
	CreateSuperUser(ctx context.Context, user *User) error
	UpdateUserToSuperUser(ctx context.Context, id int32) error
	UpdateUserToNormalUser(ctx context.Context, id int32) error
	UpdateUserInformation(ctx context.Context, user *User) error
	UpdateUserVerification(ctx context.Context, id int32) error
	DeleteUserByID(ctx context.Context, id int32) error
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

type UpdateUserInformationPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type ResendVerificationEmailPayload struct {
	Email string `json:"email" validate:"required,email"`
}
