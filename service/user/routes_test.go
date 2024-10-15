package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"

	"github.com/jayden1905/event-registration-software/types"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	app := fiber.New()

	// Define routes
	app.Get("/user", handler.handleGetUser)

	t.Run("should handle get user by token passed via cookie", func(t *testing.T) {
		// Create a new request
		req := httptest.NewRequest(fiber.MethodGet, "/user", nil)

		// Set the cookie for the user token in the request
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "valid-token",
			Path:  "/",
		})

		// Execute the test request
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error during test: %v", err)
		}

		// Check for the expected status code (200 OK)
		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode, "Status code should be 200 for valid user ID")
	})
}

// Mock implementation of the UserStore interface
type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) GetUserByID(id int32) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) CreateUser(ctx context.Context, user *types.User) error {
	return nil
}

func (m *mockUserStore) GetUserRoleByID(id int32) (string, error) {
	return "", nil
}

func (m *mockUserStore) CreateSuperUser(ctx context.Context, user *types.User) error {
	return nil
}

func (m *mockUserStore) UpdateUserToSuperUser(ctx context.Context, id int32) error {
	return nil
}

func (m *mockUserStore) UpdateUserToNormalUser(ctx context.Context, id int32) error {
	return nil
}
