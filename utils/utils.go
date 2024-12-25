package utils

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/skip2/go-qrcode"

	"github.com/jayden1905/event-registration-software/config"
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

// Function to upload an image to Cloudinary
func UploadImageToCloudinary(img []byte, folder string, format string) (string, error) {
	// Initialize Cloudinary client
	cld, err := cloudinary.NewFromParams(config.Envs.CloudinaryCloudName, config.Envs.CloudinaryAPIKey, config.Envs.CloudinarySecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	// Upload image to Cloudinary
	uploadParams := uploader.UploadParams{
		Folder: folder,
		Format: format,
	}

	// Upload the image from the byte slice
	uploadResult, err := cld.Upload.Upload(context.Background(), bytes.NewReader(img), uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %v", err)
	}

	// Return the secure URL of the uploaded image
	return uploadResult.SecureURL, nil
}

// Function to generate QR Code image and upload to Cloudinary
func GenerateQRCodeImage(data string) (string, error) {
	qrCode, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}

	// Create a new image
	img, err := qrCode.PNG(256)
	if err != nil {
		return "", err
	}

	// Upload the image to Cloudinary
	cloudinaryURL, err := UploadImageToCloudinary(img, "qr-codes", "png")
	if err != nil {
		return "", err
	}

	return cloudinaryURL, nil
}

func DeleteQrImageFromCloudinary(url string) error {
	// Initialize Cloudinary client
	cld, err := cloudinary.NewFromParams(config.Envs.CloudinaryCloudName, config.Envs.CloudinaryAPIKey, config.Envs.CloudinarySecretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	// Extract the public ID from the URL
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]

	publicID := strings.TrimSuffix(lastPart, ".png")

	folderPath := "qr-codes"
	publicID = fmt.Sprintf("%s/%s", folderPath, publicID)

	// Delete the image from Cloudinary
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %v", err)
	}

	return nil
}
