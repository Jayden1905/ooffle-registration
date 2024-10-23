package api

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jayden1905/event-registration-software/cmd/pkg/database"
	"github.com/jayden1905/event-registration-software/config"
	"github.com/jayden1905/event-registration-software/service/email"
	"github.com/jayden1905/event-registration-software/service/event"
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
		AllowOrigins:     config.Envs.PublicHosts,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization, Accept",
		AllowCredentials: true,
	}))
	// Define the api group
	api := app.Group("/api/v1")

	// Define the user store and handler
	userStore := user.NewStore(s.db)
	mailer := email.NewEmailService()
	userHandler := user.NewHandler(userStore, mailer)

	// Define the event store and handler
	eventStore := event.NewStore(s.db)
	eventHandler := event.NewHandler(eventStore, userStore)

	// Register the routes in v1 group
	userHandler.RegisterRoutes(api)
	eventHandler.RegisterRoutes(api)

	app.Use("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	app.Use("/error", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	})

	log.Println("API Server is running on: ", s.addr)
	return app.Listen(s.addr)
}
