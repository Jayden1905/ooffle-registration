package user

import (
	"fmt"

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
	router.Post("/user/login", h.handleLogin)
	router.Post("/user/register", h.handleRegister)
	router.Post("/user/super-user", h.handleCreateSuperUser)
	router.Get("/user", auth.WithJWTAuth(h.handleGetUser, h.store))
}

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
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
		HTTPOnly: true,
		Secure:   true,     // Set to true in production (HTTPS)
		SameSite: "Strict", // Prevent CSRF attacks
		Path:     "/",      // Valid for the entire site
		MaxAge:   int(config.Envs.JWTExpirationInSeconds),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "expires_in": fmt.Sprintf("%d", config.Envs.JWTExpirationInSeconds)})
}

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

func (h *Handler) handleGetUser(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	u, err := h.store.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(types.DatabaseUserToUser(u))
}

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
