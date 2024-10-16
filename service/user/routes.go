package user

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/jayden1905/event-registration-software/config"
	"github.com/jayden1905/event-registration-software/service/auth"
	"github.com/jayden1905/event-registration-software/types"
	"github.com/jayden1905/event-registration-software/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes for Fiber
func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/user/login", auth.BlockIfAuthenticated(h.handleLogin))
	router.Post("/user/logout", h.handleLogout)
	router.Post("/user/register", h.handleRegister)
	router.Patch("/user/super-user", h.handleCreateSuperUser)
	router.Put("/user/update-user/:id", auth.WithJWTAuth(h.handleUpdateUserInformation, h.store))
	router.Get("/user/current-user", auth.WithJWTAuth(h.handleGetCurrentUser, h.store))
	router.Get("/users", auth.WithJWTAuth(h.handleGetAllUsers, h.store))
	router.Get("/user/:id", auth.WithJWTAuth(h.handleGetUserByID, h.store))
	router.Delete("/user/:id", auth.WithJWTAuth(h.handleDeleteUser, h.store))
}

// Handler for registring a new user
func (h *Handler) handleRegister(c *fiber.Ctx) error {
	// Parse JSON payload
	var payload types.RegisterUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, validationErr := utils.ValidatePayload(payload)
	if validationErr != nil {
		// Return the invalid fields if validation fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user already exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("User with email %s already exists", payload.Email)})
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	// Create a new user
	err = h.store.CreateUser(c.Context(), &types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Hanlder for login
func (h *Handler) handleLogin(c *fiber.Ctx) error {
	// Parse JSON payload
	var payload types.LoginUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, err := utils.ValidatePayload(payload)
	if err != nil {
		// Return the invalid fields if validation fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user exists by email
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You don't have an account. Please register"})
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(u.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the token with HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,     // Disallow JS access to the cookie
		Secure:   true,     // Set to true in production (HTTPS)
		SameSite: "Strict", // Prevent CSRF attacks
		Path:     "/",      // Valid for the entire site
		MaxAge:   int(config.Envs.JWTExpirationInSeconds),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "expires_in": fmt.Sprintf("%d", config.Envs.JWTExpirationInSeconds)})
}

// Handler for logout
func (h *Handler) handleLogout(c *fiber.Ctx) error {
	// Clear the token cookie by setting an expired cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,     // Disallow JS access to the cookie
		Secure:   true,     // Set to true in production (HTTPS)
		SameSite: "Strict", // Prevent CSRF attacks
		Path:     "/",      // Valid for the entire site
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}

// Handler for creating a super user
func (h *Handler) handleCreateSuperUser(c *fiber.Ctx) error {
	var payload types.RegisterUserPayload

	// Parse JSON payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, err := utils.ValidatePayload(payload)
	if err != nil {
		// Return the invalid fields if validation fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	// Check if the user already exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		// compare password
		if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
		}

		// check if user is already a super user
		role, err := h.store.GetUserRoleByID(u.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
		}
		if role == "super_user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already a super user"})
		}

		// update user to super user
		h.store.UpdateUserToSuperUser(c.Context(), u.ID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"role": role,
		})
	}

	// Create a new super user
	err = h.store.CreateSuperUser(c.Context(), &types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Super user created successfully"})
}

// Handler for deleting a user
func (h *Handler) handleDeleteUser(c *fiber.Ctx) error {
	// get user id from context
	userID := auth.GetUserIDFromContext(c)

	// check if user is a super user
	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}
	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	paramsID := c.Params("id")
	// convert id to int
	intID, err := strconv.Atoi(paramsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid id: %v", err)})
	}
	// convert id to int32
	id := int32(intID)

	// Check if the user exists in the database
	_, err = h.store.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}

	// Check if the user is super_user and is trying to delete themselves
	if id == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You cannot delete yourself"})
	}

	// delete user
	if err := h.store.DeleteUserByID(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error deleting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

// Handler for updating user information
func (h *Handler) handleUpdateUserInformation(c *fiber.Ctx) error {
	var payload types.UpdateUserInformationPayload

	// Parse JSON payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, err := utils.ValidatePayload(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Get user id from context
	userID := auth.GetUserIDFromContext(c)

	// Check if the user is super_user
	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}

	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	stringID := c.Params("id")

	// convert id to int
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid id: %v", err)})
	}

	// convert id to int32
	id := int32(intID)

	// Check if the user is exists in the database
	_, err = h.store.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}

	// Update the user information
	if err := h.store.UpdateUserInformation(c.Context(), &types.User{
		ID:        id,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error updating user information: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User information updated successfully"})
}

// Handler for getting the current user
func (h *Handler) handleGetCurrentUser(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	u, err := h.store.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

// Handler for getting all users
func (h *Handler) handleGetAllUsers(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// check if user is a super user
	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}
	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	users, err := h.store.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting all users: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// Handler for getting a user by id
func (h *Handler) handleGetUserByID(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}
	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	stringID := c.Params("id")

	// convert id to int
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid id: %v", err)})
	}

	// convert id to int32
	id := int32(intID)

	u, err := h.store.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(u)
}
