package handlers

import (
	"Smd/services"
	"Smd/types"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type MockFileService struct {
	services.FileService
	LastUploadedFile types.File
}

func (m *MockFileService) UploadFile(f types.File) error {
	m.LastUploadedFile = f
	return nil
}

func TestUploadFileHandler(t *testing.T) {
	// Create a temporary file to upload
	tempFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_, err = tempFile.WriteString("Test file content")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Prepare the multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadFile", filepath.Base(tempFile.Name()))
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	part.Write(fileContent)
	writer.Close()

	// Create the request
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Mock the FileService
	mockService := &MockFileService{}

	// Call the handler
	rr := httptest.NewRecorder()
	UploadFileHandler(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Assert the file metadata in the mock service
	if mockService.LastUploadedFile.Name != filepath.Base(tempFile.Name()) {
		t.Errorf("Uploaded file name mismatch: got %v want %v", mockService.LastUploadedFile.Name, filepath.Base(tempFile.Name()))
	}

	// Add more assertions as needed
}
