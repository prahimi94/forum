package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Allowed file extensions
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// Check if file extension is allowed
func IsAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}

func FileUpload(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Check file extension before proceeding
	if !IsAllowedExtension(handler.Filename) {
		return "", fmt.Errorf("file type not allowed: %s", handler.Filename)
	}

	uploadDir := "static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm) // Ensure directory exists

	fileUUID, err := GenerateUuid()
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(handler.Filename)
	filePath := filepath.Join(uploadDir, fileUUID+fileExt)
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return "", err
	}
	defer outFile.Close()

	// Copy the uploaded file to the new location
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Println("Error saving file:", err)
		return "", err
	}

	return fileUUID + fileExt, nil
}
