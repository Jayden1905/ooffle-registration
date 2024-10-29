package attendee

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/jayden1905/event-registration-software/service/auth"
	"github.com/jayden1905/event-registration-software/types"
)

type Handler struct {
	store      types.AttendeeStore
	eventStore types.EventStore
	userStore  types.UserStore
}

func NewHandler(store types.AttendeeStore, eventStore types.EventStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, eventStore: eventStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/events/:event_id/attendees", auth.WithJWTAuth(h.handleGetAttendeesPaginated, h.userStore))
}

// Handler to get all the attendees paginated from database
func (h *Handler) handleGetAttendeesPaginated(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// check if the user is the owner of the event
	eventIDString := c.Params("event_id")
	eventID, err := strconv.Atoi(eventIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	event, err := h.eventStore.GetEventByID(int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get event",
		})
	}

	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	const (
		defaultPageSize = 10
		maxPageSize     = 100
	)

	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page := 1
	pageSize := defaultPageSize

	// Parse page if provided
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse pageSize if provided
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= maxPageSize {
			pageSize = ps
		}
	}

	attendees, err := h.store.GetAllAttendeesPaginated(int32(page), int32(pageSize))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendees",
		})
	}

	return c.Status(fiber.StatusOK).JSON(attendees)
}
