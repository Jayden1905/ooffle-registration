package user

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"

	"github.com/jayden1905/event-registration-software/cmd/pkg/database"
)

type Store struct {
	db *database.Queries
}

// NewStore initializes the Store with the database queries
func NewStore(db *database.Queries) *Store {
	return &Store{db: db}
}

// GetUserByEmail fetches a user by email from the database
func (s *Store) GetUserByEmail(email string) (*database.User, error) {
	user, err := s.db.GetUserByEmail(context.Background(), email) // Use the SQLC-generated method
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID fetches a user by ID from the database
func (s *Store) GetUserByID(id int32) (*database.User, error) {
	user, err := s.db.GetUserByID(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserRoleByID fetches the role of a user by ID from the database
func (s *Store) GetUserRoleByID(id int32) (string, error) {
	role, err := s.db.GetUserRoleByUserID(context.Background(), id)
	if err != nil {
		return "", err
	}

	stringRole := string(role)

	return stringRole, nil
}

// CreateUser creates a new user in the database
func (s *Store) CreateUser(ctx context.Context, user *database.User) error {
	err := s.db.CreateNormalUser(ctx, database.CreateNormalUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateSuperUser creates a new super user in the database
func (s *Store) CreateSuperUser(ctx context.Context, user *database.User) error {
	err := s.db.CreateSuperUser(ctx, database.CreateSuperUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		return err
	}

	return nil
}

// Update the user to a super user
func (s *Store) UpdateUserToSuperUser(ctx context.Context, id int32) error {
	err := s.db.UpdateUserToSuperUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Update the user to a normal user
func (s *Store) UpdateUserToNormalUser(ctx context.Context, id int32) error {
	err := s.db.UpdateUserToNormalUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
