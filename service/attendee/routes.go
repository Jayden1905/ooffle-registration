package attendee

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/jayden1905/event-registration-software/service/auth"
	"github.com/jayden1905/event-registration-software/types"
	"github.com/jayden1905/event-registration-software/utils"
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
	router.Get("/event/:event_id/attendees", auth.WithJWTAuth(h.handleGetAttendeesPaginated, h.userStore))
	router.Post("/event/add_attendee", auth.WithJWTAuth(h.handleCreateNewAttendee, h.userStore))
	router.Delete("/event/:event_id/attendees/:attendee_id", auth.WithJWTAuth(h.handleDeleteAttendeeByID, h.userStore))
	router.Delete("/event/:event_id/attendees", auth.WithJWTAuth(h.handleDeleteAllAttendeesByEventID, h.userStore))
	router.Patch("/event/attendees/:attendee_id", auth.WithJWTAuth(h.handleUpdateAttendeeByID, h.userStore))
	router.Post("/event/:event_id/attendees/import", auth.WithJWTAuth(h.handleImportAttendeesFromCSV, h.userStore))
}

// Handler to create a new attendee
func (h *Handler) handleCreateNewAttendee(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	var payload types.CreateAttendeePayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	// Validate the payload
	if invalidFields, err := utils.ValidatePayload(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid payload",
			"fields": invalidFields,
		})
	}

	// Check if the user is the owner of the event
	event, err := h.eventStore.GetEventByID(payload.EventID)
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

	// Check if the attendee with same email already exists
	atte, err := h.store.GetAttendeeByEmail(payload.Email)
	if atte != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Attendee with same email already exists",
		})
	}

	// Generate QR code
	qrCode, err := utils.GenerateQRCodeBase64(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate QR code",
		})
	}

	attendee := &types.Attendee{
		FristName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		EventID:     payload.EventID,
		QrCode:      qrCode,
		CompanyName: payload.CompanyName,
		Title:       payload.Title,
		TableNo:     payload.TableNo,
		Role:        payload.Role,
		Attendance:  false,
	}

	if err := h.store.CreateAttendee(c.Context(), attendee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create attendee",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(attendee)
}

// Handler to import attendees from a CSV file
func (h *Handler) handleImportAttendeesFromCSV(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// Check if the user is the owner of the event
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

	// Parse the file from the request
	file, err := c.FormFile("import")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file",
		})
	}
	// Check if the file is a CSV file
	if file.Header.Get("Content-Type") != "text/csv" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type",
		})
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}

	// Parse the CSV file
	csvReader := csv.NewReader(fileContent)
	records, err := csvReader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse CSV file",
		})
	}

	// Skip the header row and insert each row as an attendee
	for _, record := range records[1:] {
		qrCodeBase64, err := utils.GenerateQRCodeBase64(record[2])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate QR code",
			})
		}

		error := h.store.CreateAttendee(c.Context(), &types.Attendee{
			FristName:   record[0],
			LastName:    record[1],
			Email:       record[2],
			EventID:     int32(eventID),
			QrCode:      qrCodeBase64,
			CompanyName: record[3],
			Title:       record[4],
			TableNo:     utils.ParseTableNo(record[5]),
			Role:        record[6],
		})
		if error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create attendee",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Attendees imported successfully",
	})
}

// Handler to update an attendee by ID
func (h *Handler) handleUpdateAttendeeByID(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// Convert the attendee ID to integer from params
	attendeeIDString := c.Params("attendee_id")
	attendeeID, err := strconv.Atoi(attendeeIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attendee ID",
		})
	}

	// Check if the attendee exists
	attendee, err := h.store.GetAttendeeByID(int32(attendeeID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Attendee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendee",
		})
	}

	// Check if the user is the owner of the event
	event, err := h.eventStore.GetEventByID(attendee.EventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Event not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get event",
		})
	}

	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var payload types.UpdateAttendeePayload
	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	// Validate the payload
	if invalidFields, err := utils.ValidatePayload(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid payload",
			"fields": invalidFields,
		})
	}

	if payload.Email != "" && payload.Email != attendee.Email {
		log.Println(payload.Email)
		// Generate QR code
		qrCode, err := utils.GenerateQRCodeBase64(payload.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate QR code",
			})
		}
		// Update the attendee by ID
		if err := h.store.UpdateAttendeeByID(int32(attendeeID), &types.Attendee{
			FristName:   payload.FirstName,
			LastName:    payload.LastName,
			Email:       payload.Email,
			QrCode:      qrCode,
			CompanyName: payload.CompanyName,
			Title:       payload.Title,
			TableNo:     payload.TableNo,
			Role:        payload.Role,
			Attendance:  payload.Attendance,
		}); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update attendee",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Attendee updated successfully with new qrcode",
		})
	}

	// Update the attendee by ID
	if err := h.store.UpdateAttendeeByID(int32(attendeeID), &types.Attendee{
		FristName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		CompanyName: payload.CompanyName,
		Title:       payload.Title,
		TableNo:     payload.TableNo,
		Role:        payload.Role,
		Attendance:  payload.Attendance,
	}); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update attendee",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendee updated successfully",
	})
}

// Handler for deleting an attendee by ID
func (h *Handler) handleDeleteAttendeeByID(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	eventIDString := c.Params("event_id")

	// Convert the event ID to integer
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

	// Check if the user is the owner of the event
	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Convert the attendee ID to integer from params
	attendeeIDString := c.Params("attendee_id")
	attendeeID, err := strconv.Atoi(attendeeIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attendee ID",
		})
	}

	// Check if the attendee exists
	_, err = h.store.GetAttendeeByID(int32(attendeeID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Attendee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendee",
		})
	}

	// Delete the attendee by ID
	if err := h.store.DeleteAttendeeByID(int32(attendeeID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete attendee",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendee deleted successfully",
	})
}

// Handler for deleting all attendees by event ID
func (h *Handler) handleDeleteAllAttendeesByEventID(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)
	eventIDString := c.Params("event_id")

	// Convert the event ID to integer
	eventID, err := strconv.Atoi(eventIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	// Retrieve the event by its ID
	event, err := h.eventStore.GetEventByID(int32(eventID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Event not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve event",
		})
	}

	// Check if the user is the owner of the event
	if event.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Delete all attendees by event ID
	if err := h.store.DeleteAllAttendeesByEventID(int32(eventID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete attendees",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendees deleted successfully",
	})
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

	attendees, err := h.store.GetAllAttendeesPaginated(int32(page), int32(pageSize), int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendees",
		})
	}

	return c.Status(fiber.StatusOK).JSON(attendees)
}
