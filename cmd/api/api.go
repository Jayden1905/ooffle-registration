package api

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jayden1905/event-registration-software/cmd/pkg/database"
	"github.com/jayden1905/event-registration-software/config"
	"github.com/jayden1905/event-registration-software/service/user"
)

type apiConfig struct {
	addr string
	db   *database.Queries
}

func NewAPIServer(addr string, db *sql.DB) *apiConfig {
	return &apiConfig{
		addr: addr,
		db:   database.New(db),
	}
}

func (s *apiConfig) Run() error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Envs.PublicHost,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// Define the user store and handler
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	// Define the api group
	api := app.Group("/api/v1")

	// Register the user routes
	userHandler.RegisterRoutes(api)

	// create a health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	// curl -X GET http://localhost:8080/health

	log.Println("API Server is running on: ", s.addr)
	return app.Listen(s.addr)
}
