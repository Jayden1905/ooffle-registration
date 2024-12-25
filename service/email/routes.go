package email

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jayden1905/event-registration-software/service/auth"
	"github.com/jayden1905/event-registration-software/types"
	"github.com/jayden1905/event-registration-software/utils"
)

type Handler struct {
	store      types.EmailTempalteStore
	eventStore types.EventStore
	userStore  types.UserStore
}

func NewHandler(store types.EmailTempalteStore, eventStore types.EventStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, eventStore: eventStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/email_templates/:event_id", auth.WithJWTAuth(h.handleGetEmailTempalteByID, h.userStore))
	router.Post("/email_templates", auth.WithJWTAuth(h.handleCreateEmailTemplate, h.userStore))
	router.Put("/email_templates", auth.WithJWTAuth(h.handleUpdateEmailTemplate, h.userStore))
}

// Handler for getting an email template by its ID
func (h *Handler) handleGetEmailTempalteByID(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// check if the user is the owner of the event
	eventIDString := c.Params("event_id")
	eventID, err := strconv.Atoi(eventIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	event, err := h.eventStore.GetEventByID(int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	emailTemplate, err := h.store.GetEmailTemplateByEventID(c.Context(), int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get email template"})
	}

	return c.Status(fiber.StatusOK).JSON(emailTemplate)
}

// Handler for creating an email template
func (h *Handler) handleCreateEmailTemplate(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	var payload types.CreateEmailTemplatePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	if invalidFields, err := utils.ValidatePayload(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid payload fields",
			"fields": invalidFields,
		})
	}

	event, err := h.eventStore.GetEventByID(payload.EventID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	emailTemplate := &types.EmailTemplate{
		EventID:     payload.EventID,
		HeaderImage: payload.HeaderImage,
		Content:     payload.Content,
		FooterImage: payload.FooterImage,
		Subject:     payload.Subject,
		BgColor:     payload.BgColor,
		Message:     payload.Message,
	}

	// check if the email template already exists
	if _, err := h.store.GetEmailTemplateByEventID(c.Context(), payload.EventID); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email template already exists"})
	}

	if err := h.store.CreateEmailTemplate(c.Context(), emailTemplate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create email template"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Email template created successfully"})
}

// Handler for updating an email template
func (h *Handler) handleUpdateEmailTemplate(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	var payload types.UpdateEmailTemplatePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	if invalidFields, err := utils.ValidatePayload(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid payload fields",
			"fields": invalidFields,
		})
	}

	event, err := h.eventStore.GetEventByID(payload.EventID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	emailTemplate := &types.EmailTemplate{
		ID:          payload.ID,
		EventID:     payload.EventID,
		HeaderImage: payload.HeaderImage,
		Content:     payload.Content,
		FooterImage: payload.FooterImage,
		Subject:     payload.Subject,
		BgColor:     payload.BgColor,
		Message:     payload.Message,
	}

	if err := h.store.UpdateEmailTemplate(c.Context(), emailTemplate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update email template"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email template updated successfully"})
}
