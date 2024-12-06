package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/skip2/go-qrcode"

	"github.com/jayden1905/event-registration-software/types"
)

var Validate = validator.New()

func ValidatePayload(payload interface{}) (map[string]string, error) {
	err := Validate.Struct(payload)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			invalidFields := make(map[string]string)
			for _, e := range validationErrors {
				invalidFields[e.Field()] = fmt.Sprintf("Validation failed on the '%s' tag", e.Tag())
			}
			return invalidFields, fmt.Errorf("invalid payload")
		}
	}
	return nil, nil
}

func IsSuperUser(userID int32, store types.UserStore) (bool, error) {
	// Get the user role from the store
	role, err := store.GetUserRoleByID(userID)
	if err != nil {
		return false, fmt.Errorf("error getting user role by id: %v", err)
	}

	// Check if the role is "super_user"
	if role == "super_user" {
		return true, nil
	}

	return false, nil
}

// Function to generate QR code as a base64 string
func GenerateQRCodeBase64(data string) (string, error) {
	qrCode, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = qrCode.Write(256, &buffer)
	if err != nil {
		return "", err
	}

	// Convert to base64 string
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

// Parsing an integer for table number
func ParseTableNo(value string) int32 {
	// remove the whitespaces
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	// convert the string to an integer
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return int32(num)
}
