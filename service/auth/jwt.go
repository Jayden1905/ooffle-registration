package auth

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/jayden1905/event-registration-software/config"
	"github.com/jayden1905/event-registration-software/types"
)

type contextKey string

const UserKey contextKey = "userID"

// CreateJWT generates a new JWT token with the given secret and userID.
func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// WithJWTAuth is a middleware for Fiber that validates the JWT token.
func WithJWTAuth(handlerFunc fiber.Handler, store types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from cookies or Authorization header
		tokenString, err := getTokenFromCookie(c)
		if err != nil {
			tokenString = c.Get("Authorization")
			if tokenString == "" {
				return permissionDenied(c)
			}
		}

		// Validate the JWT token
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			return permissionDenied(c)
		}
		if !token.Valid {
			log.Printf("token is invalid")
			return permissionDenied(c)
		}

		// Extract the userID from JWT claims
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			return permissionDenied(c)
		}

		// Fetch the user from the database
		u, err := store.GetUserByID(int32(userID))
		if err != nil {
			log.Printf("error getting user by id: %v", err)
			return permissionDenied(c)
		}

		// Set userID in context (using Fiber's Locals)
		c.Locals(UserKey, u.ID)

		// Call the next handler
		return handlerFunc(c)
	}
}

// Helper function to extract JWT token from cookie in Fiber
func getTokenFromCookie(c *fiber.Ctx) (string, error) {
	token := c.Cookies("token")
	if token == "" {
		return "", fmt.Errorf("token not found in cookies")
	}
	return token, nil
}

// Helper function to validate a JWT token
func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

// Helper function to send a permission denied response in Fiber
func permissionDenied(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
}

func GetUserIDFromContext(c *fiber.Ctx) int32 {
	userID, ok := c.Locals(UserKey).(int32)
	if !ok {
		return 0
	}
	return userID
}
