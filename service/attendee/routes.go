package attendee

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"

	"github.com/jayden1905/event-registration-software/service/auth"
	"github.com/jayden1905/event-registration-software/service/email"
	"github.com/jayden1905/event-registration-software/types"
	"github.com/jayden1905/event-registration-software/utils"
)

type Handler struct {
	store      types.AttendeeStore
	eventStore types.EventStore
	userStore  types.UserStore
	emailStore types.EmailTempalteStore
	mailer     email.Mailer
}

func NewHandler(store types.AttendeeStore, eventStore types.EventStore, userStore types.UserStore, emailStore types.EmailTempalteStore, mailer email.Mailer) *Handler {
	return &Handler{store: store, eventStore: eventStore, userStore: userStore, emailStore: emailStore, mailer: mailer}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/event/:event_id/attendees", auth.WithJWTAuth(h.handleGetAttendeesPaginated, h.userStore))
	router.Get("/event/:event_id/attendees/count", auth.WithJWTAuth(h.handleGetAttendeesRowCount, h.userStore))
	router.Get("/attendees/:attendee_id", auth.WithJWTAuth(h.handleGetAttendeeByID, h.userStore))
	router.Get("/event/:event_id/attendees/all", auth.WithJWTAuth(h.handleGetAllAttendees, h.userStore))
	router.Post("/event/add_attendee", auth.WithJWTAuth(h.handleCreateNewAttendee, h.userStore))
	router.Delete("/event/:event_id/attendees/:attendee_id", auth.WithJWTAuth(h.handleDeleteAttendeeByID, h.userStore))
	router.Delete("/event/:event_id/attendees", auth.WithJWTAuth(h.handleDeleteAllAttendeesByEventID, h.userStore))
	router.Put("/event/attendees/:attendee_id", auth.WithJWTAuth(h.handleUpdateAttendeeByID, h.userStore))
	router.Post("/event/:event_id/attendees/import", auth.WithJWTAuth(h.handleImportAttendeesFromCSV, h.userStore))
	router.Post("/event/:event_id/attendees/send_invitation", auth.WithJWTAuth(h.handleSendInvitationEmail, h.userStore))
	router.Post("/attendees/mark_attendance/:attendee_email", auth.WithJWTAuth(h.handleMarkAttendeeAttendance, h.userStore))
}

func (h *Handler) handleGetAttendeeByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(attendee)
}

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
			"error":  "Invalid payload fields",
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendee",
		})
	}

	// Generate and upload QR code concurrently with error handling
	qrCodeChannel := make(chan string, 1)
	errorChannel := make(chan error, 1)

	go func() {
		// Generate the QR code image
		qrCode, err := utils.GenerateQRCodeImage(payload.Email)
		if err != nil {
			errorChannel <- err // Send error if QR code generation fails
			return
		}

		qrCodeChannel <- qrCode
	}()

	// Wait for either the QR code URL or an error from the channels
	select {
	case qrCodeURL := <-qrCodeChannel:
		// Proceed if QR code upload was successful
		attendee := &types.Attendee{
			FirstName:   payload.FirstName,
			LastName:    payload.LastName,
			Email:       payload.Email,
			EventID:     payload.EventID,
			QrCode:      qrCodeURL,
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

	case err := <-errorChannel:
		// If an error occurred during QR code generation or upload
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate or upload QR code: %v", err),
		})
	}
}

// Error result type for attendee errors
type errorAttendeeResult struct {
	attendee types.Attendee
	err      error
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
	errorAttendees := []types.Attendee{}
	var wg sync.WaitGroup
	attendeeErrorsChan := make(chan errorAttendeeResult, len(records)-1) // Buffered channel to hold error results

	for _, record := range records[1:] {
		wg.Add(1) // Increment WaitGroup counter for each goroutine

		go func(record []string) {
			defer wg.Done() // Decrement the counter when the goroutine finishes

			// Generate QR code
			qrCode, err := utils.GenerateQRCodeImage(record[2])
			if err != nil {
				attendeeErrorsChan <- errorAttendeeResult{attendee: types.Attendee{}, err: fmt.Errorf("failed to generate QR code: %v", err)}
				return
			}

			// Check if the attendee with the same email already exists in the same event
			attendeeExist, err := h.store.GetAttendeeByEmail(record[2])
			if attendeeExist != nil && attendeeExist.EventID == int32(eventID) {
				attendeeErrorsChan <- errorAttendeeResult{attendee: *attendeeExist, err: nil}
				return
			}

			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				attendeeErrorsChan <- errorAttendeeResult{attendee: types.Attendee{}, err: fmt.Errorf("failed to get attendee: %v", err)}
				return
			}

			// Create the attendee
			attendee := &types.Attendee{
				FirstName:   record[0],
				LastName:    record[1],
				Email:       record[2],
				EventID:     int32(eventID),
				QrCode:      qrCode,
				CompanyName: record[3],
				Title:       record[4],
				TableNo:     utils.ParseTableNo(record[5]),
				Role:        record[6],
				Attendance:  false,
			}

			// Insert attendee
			if err := h.store.CreateAttendee(c.Context(), attendee); err != nil {
				attendeeErrorsChan <- errorAttendeeResult{attendee: *attendee, err: fmt.Errorf("failed to create attendee: %v", err)}
				return
			}

			// Success: No error for this attendee
			attendeeErrorsChan <- errorAttendeeResult{attendee: *attendee, err: nil}
		}(record)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Close the channel after all goroutines are done
	close(attendeeErrorsChan)

	// Collect errors if any
	for result := range attendeeErrorsChan {
		if result.err != nil {
			errorAttendees = append(errorAttendees, result.attendee)
		}
	}

	// Return a partial success message if there are any errors
	if len(errorAttendees) > 0 {
		return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
			"message":   "Attendees imported with some errors",
			"error":     "Some attendees already exist or failed to create",
			"attendees": errorAttendees,
		})
	}

	// Successful import
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
		qrCode, err := utils.GenerateQRCodeImage(payload.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate QR code",
			})
		}
		// Update the attendee by ID
		if err := h.store.UpdateAttendeeByID(int32(attendeeID), &types.Attendee{
			FirstName:   payload.FirstName,
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
		FirstName:   payload.FirstName,
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

	// Delete qr image from Cloudinary
	err = utils.DeleteQrImageFromCloudinary(attendee.QrCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete image",
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

	// Get all attendees by event ID
	attendees, err := h.store.GetAllAttendees(int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendees",
		})
	}

	// Channel to capture image deletion errors
	errChan := make(chan error, len(attendees))

	// Delete qr images from Cloudinary concurrently
	for _, attendee := range attendees {
		go func(attendee types.Attendee) {
			if deleteErr := utils.DeleteQrImageFromCloudinary(attendee.QrCode); deleteErr != nil {
				errChan <- fmt.Errorf("failed to delete image for attendee %d: %v", attendee.ID, deleteErr)
			} else {
				errChan <- nil
			}
		}(*attendee)
	}

	// Wait for all image deletion results
	var deletionErrors []string
	for i := 0; i < len(attendees); i++ {
		if deleteErr := <-errChan; deleteErr != nil {
			deletionErrors = append(deletionErrors, deleteErr.Error())
		}
	}

	// If there were any errors deleting images, log and proceed
	if len(deletionErrors) > 0 {
		log.Printf("Image deletion errors: %v", deletionErrors)
	}

	// Delete all attendees by event ID
	if err := h.store.DeleteAllAttendeesByEventID(int32(eventID)); err != nil {
		log.Printf("Error deleting attendees for event ID: %d: %v", eventID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete attendees",
		})
	}

	// Provide a response with success
	responseMessage := "Attendees deleted successfully"
	if len(deletionErrors) > 0 {
		// Include a message about image deletion issues
		responseMessage = fmt.Sprintf("%s, but some images were not deleted: %v", responseMessage, deletionErrors)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": responseMessage,
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

// Handler to get all the attendees from database
func (h *Handler) handleGetAllAttendees(c *fiber.Ctx) error {
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

	attendees, err := h.store.GetAllAttendees(int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendees",
		})
	}

	return c.Status(fiber.StatusOK).JSON(attendees)
}

// Handler to get attendees row count from database
func (h *Handler) handleGetAttendeesRowCount(c *fiber.Ctx) error {
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

	rowCount, err := h.store.GetAttendeeRowCount(int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendees row count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"row_count": rowCount,
	})
}

// handler to send invitation email to attendees
func (h *Handler) handleSendInvitationEmail(c *fiber.Ctx) error {
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

	// Get the email template by event ID
	emailTemplate, err := h.emailStore.GetEmailTemplateByEventID(c.Context(), int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get email template",
		})
	}

	// Get all the attendees by event ID
	attendees, err := h.store.GetAllAttendees(int32(eventID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendees",
		})
	}

	emailTmp := &types.EmailTemplate{
		ID:          emailTemplate.ID,
		EventID:     emailTemplate.EventID,
		HeaderImage: emailTemplate.HeaderImage,
		Content:     emailTemplate.Content,
		FooterImage: emailTemplate.FooterImage,
		Subject:     emailTemplate.Subject,
		BgColor:     emailTemplate.BgColor,
		Message:     emailTemplate.Message,
	}

	// Use a wait group to synchronize goroutines
	var wg sync.WaitGroup
	// Create a channel to handle errors from goroutines
	errorChannel := make(chan error, len(attendees))

	// Limit the number of concurrent goroutines
	concurrencyLimit := 10
	semaphore := make(chan struct{}, concurrencyLimit)

	// Send invitation email to each attendee concurrently
	for _, attendee := range attendees {
		// Increment wait group counter
		wg.Add(1)

		// Acquire a spot in the semaphore
		semaphore <- struct{}{}

		go func(att *types.Attendee) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the spot in the semaphore

			// Send email
			err := h.mailer.SendInvitationEmail(att, emailTmp)
			if err != nil {
				log.Println("Error sending email to:", att.Email, err)
				errorChannel <- err
				return
			}
		}(attendee)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChannel)

	// If there were any errors in sending emails
	if len(errorChannel) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Some emails failed to send",
		})
	}

	// Successfully sent all emails
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Invitation emails sent successfully",
	})
}

func (h *Handler) handleMarkAttendeeAttendance(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// Convert the attendee ID to integer from params
	attendeeEmail := c.Params("attendee_email")

	// Check if the attendee exists
	attendee, err := h.store.GetAttendeeByEmail(attendeeEmail)
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

	// Check if the attendee has already been marked
	if attendee.Attendance {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Attendance already marked",
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

	// Mark the attendance of the attendee
	if err := h.store.UpdateAttendeeByID(attendee.ID, &types.Attendee{
		FirstName:   attendee.FirstName,
		LastName:    attendee.LastName,
		Email:       attendee.Email,
		EventID:     attendee.EventID,
		QrCode:      attendee.QrCode,
		CompanyName: attendee.CompanyName,
		Title:       attendee.Title,
		TableNo:     attendee.TableNo,
		Role:        attendee.Role,
		Attendance:  true,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark attendance",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendance marked successfully",
	})
}
