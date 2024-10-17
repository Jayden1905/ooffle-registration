package event

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/jayden1905/event-registration-software/service/auth"
	"github.com/jayden1905/event-registration-software/types"
	"github.com/jayden1905/event-registration-software/utils"
)

type Handler struct {
	store     types.EventStore
	userStore types.UserStore
}

func NewHandler(store types.EventStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/events", auth.WithJWTAuth(h.handleGetAllEvents, h.userStore))
	router.Post("/event/create", auth.WithJWTAuth(h.handleCreateEvent, h.userStore))
	router.Put("/event/update/:id", auth.WithJWTAuth(h.handleUpdateEvent, h.userStore))
	router.Delete("/event/delete/:id", auth.WithJWTAuth(h.handleDeleteEventByEventID, h.userStore))
	// router.Delete("/events/delete", auth.WithJWTAuth(h.handleDeleteAllEvents, h.userStore))
}

// handleGetAllEvents fetches all events from the database
func (h *Handler) handleGetAllEvents(c *fiber.Ctx) error {
	// get user id from the context
	userID := auth.GetUserIDFromContext(c)

	events, err := h.store.GetAllEvents(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(events)
}

// handleCreateEvent creates a new event in the database
func (h *Handler) handleCreateEvent(c *fiber.Ctx) error {
	// Parse the request payload
	var payload types.CreateEventPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Println(err)
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

	// Check if the event with the same title already exists
	if _, err := h.store.GetEventByTitle(payload.Title); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Event with the same title already exists"})
	}

	// get user id from the context
	userID := auth.GetUserIDFromContext(c)

	// Create a new event
	if err := h.store.CreateNewEvent(c.Context(), &types.Event{
		Title:       payload.Title,
		Description: payload.Description,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		Location:    payload.Location,
		UserID:      userID,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(payload)
}

// handleUpdateEvent updates an existing event in the database
func (h *Handler) handleUpdateEvent(c *fiber.Ctx) error {
	// Parse the request payload
	var payload types.CreateEventPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Println(err)
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

	// get user id from the context
	userID := auth.GetUserIDFromContext(c)

	// get event id from the context
	eventIDString := c.Params("id")
	eventIDInt, err := strconv.Atoi(eventIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event id"})
	}
	eventID := int32(eventIDInt)

	// Check if the event exists
	event, err := h.store.GetEventByID(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Event not found"})
	}

	// check if the user is the owner of the event
	if userID != event.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to update this event"})
	}

	// Update the event
	if err := h.store.UpdateEvent(c.Context(), &types.Event{
		EventID:     eventID,
		Title:       payload.Title,
		Description: payload.Description,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		Location:    payload.Location,
		UserID:      userID,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Event updated successfully"})
}

// handleDeleteEventByEventID deletes an event from the database
func (h *Handler) handleDeleteEventByEventID(c *fiber.Ctx) error {
	// get user id from the context
	userID := auth.GetUserIDFromContext(c)

	// get event id from the context
	eventIDString := c.Params("id")
	eventIDInt, err := strconv.Atoi(eventIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event id"})
	}
	eventID := int32(eventIDInt)

	// Check if the event exists
	event, err := h.store.GetEventByID(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Event not found"})
	}

	// check if the user is the owner of the event
	if userID != event.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to delete this event"})
	}

	// Delete the event
	if err := h.store.DeleteEvent(c.Context(), eventID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Event deleted successfully"})
}

// handle DeleteAllEvents deletes all events from the database
func (h *Handler) handleDeleteAllEvents(c *fiber.Ctx) error {
	// get user id from the context
	userID := auth.GetUserIDFromContext(c)

	err := h.store.DeleteAllEvents(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "All events deleted successfully"})
}
