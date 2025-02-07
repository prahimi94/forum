package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func FileUpload(file multipart.File, handler *multipart.FileHeader) (string, error) {
	uploadDir := "static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm) // Ensure directory exists

	filePath := filepath.Join(uploadDir, handler.Filename)
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
	fmt.Printf("File uploaded successfully: %s\n", filePath)
	return filePath, nil
}
