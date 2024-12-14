package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

// CreateVerificationToken generates a new JWT token with the given secret and userID.
func GenerateVerificationToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(5 * time.Minute).Unix(), // Token expires in 5 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(config.Envs.JWTSecret)

	return token.SignedString(secret)
}

// WithJWTAuth is a middleware for Fiber that validates the JWT token.
func WithJWTAuth(handlerFunc fiber.Handler, store types.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from cookies or Authorization header
		tokenString, err := getTokenFromCookie(c)
		if err != nil || tokenString == "" {
			tokenString = c.Get("Authorization")
		}

		if tokenString == "" {
			return permissionDenied(c)
		}

		// Validate the JWT token
		token, err := ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
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

// BlockIfAuthenticated is a middleware for Fiber that blocks the request if the user is authenticated.
func BlockIfAuthenticated(handlerFunc fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from cookies or Authorization header
		tokenString, err := getTokenFromCookie(c)
		if err != nil || tokenString == "" {
			tokenString = c.Get("Authorization")
		}

		if tokenString == "" {
			return handlerFunc(c)
		}

		// Validate the JWT token
		token, err := ValidateToken(tokenString)

		// If the token is valid, block the request
		if err == nil && token.Valid {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "User is already authenticated"})
		}

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

// Helper function to validate the verification token
func ValidateVerificationToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
	// Check for any errors during parsing
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// Check if the error was due to token expiration
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", fmt.Errorf("token has expired")
			} else {
				return "", fmt.Errorf("token is invalid: %v", err)
			}
		}
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error parsing claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("error parsing email")
	}

	return email, nil
}

// Helper function to validate a JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
	// Check for any errors during parsing
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// Check if the error was due to token expiration
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token has expired")
			} else {
				return nil, fmt.Errorf("token is invalid: %v", err)
			}
		}
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	return token, nil
}

// Helper function to send a permission denied response in Fiber
func permissionDenied(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Permission denied"})
}

// GetUserIDFromContext extracts the userID from Fiber's context
func GetUserIDFromContext(c *fiber.Ctx) int32 {
	userID, ok := c.Locals(UserKey).(int32)
	if !ok {
		return 0
	}
	return userID
}

// hashUserID hashes the userID using SHA-256
func hashUserID(userID int32) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d", userID)))
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateRateLimiter returns a Fiber middleware for rate limiting
func CreateRateLimiter(maxRequests int, expiration time.Duration, customErrMessage string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use userID from Locals (set by authentication middleware)
			userID := GetUserIDFromContext(c)
			if userID == 0 {
				// If user is not authenticated (no userID), rate limit by IP
				return c.IP()
			}
			return fmt.Sprintf("user:%v", hashUserID(userID))
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Custom error response when limit is reached
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": customErrMessage,
			})
		},
	})
}
