package services

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseAndValidateFile(t *testing.T) {
	fs := &fileService{}

	tests := []struct {
		name    string
		errMsg  string
		body    []byte
		wantErr bool
	}{
		{
			name:    "FileTooBig",
			body:    bytes.Repeat([]byte("a"), 10<<20+1), // More than 10MB
			wantErr: true,
			errMsg:  "file too big",
		},
		{
			name:    "Success",
			body:    []byte("test data"), // Less than 10MB
			wantErr: false,
		},
		{
			name:    "EmptyBody",
			body:    []byte(""),
			wantErr: false,
		},
		{
			name:    "NilBody",
			body:    nil,
			wantErr: true,
			errMsg:  "error retrieving the file",
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to hold the request body
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Only create a file part if the body is not empty or nil
			if tt.body != nil {
				// Create a form file writer for the file
				part, err := writer.CreateFormFile("uploadFile", "testfile")
				if err != nil {
					t.Fatal(err)
				}

				// Copy the test data into the form file writer
				_, err = part.Write(tt.body)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Close the multipart writer to finish writing the request body
			err := writer.Close()
			if err != nil {
				t.Fatal(err)
			}

			// Create a new request with the body
			r, err := http.NewRequest("POST", "/upload", body)
			if err != nil {
				t.Fatal(err)
			}

			// Set the content type of the request to multipart/form-data
			if writer != nil {
				r.Header.Set("Content-Type", writer.FormDataContentType())
			}

			w := httptest.NewRecorder()

			_, _, err = fs.ParseAndValidateFile(r, w)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("expected '%s' error, got %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
func TestSaveFile(t *testing.T) {
	fs := &fileService{}

	tests := []struct {
		name        string
		file        multipart.File
		destination string
		wantErr     bool
	}{
		{
			name: "Success",
			file: func() multipart.File {
				file, err := os.CreateTemp("", "test")
				if err != nil {
					t.Fatal(err)
				}
				file.WriteString("test data")
				file.Seek(0, 0)
				return file
			}(), // Simulate a file with "test data"
			destination: func() string {
				dir, err := os.MkdirTemp("", "destination")
				if err != nil {
					t.Fatal(err)
				}
				return dir
			}(),
			wantErr: false,
		},
		{
			name:        "NilFile",
			file:        nil,
			destination: "",
			wantErr:     true,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.SaveFile(tt.file, tt.name, tt.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestDeleteFile(t *testing.T) {
// 	fs := &fileService{}

// 	tests := []struct {
// 		name    string
// 		file    string // Assume DeleteFile takes a filename
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Success",
// 			file:    "testfile", // Assume this file exists
// 			wantErr: false,
// 		},
// 		{
// 			name:    "NonExistentFile",
// 			file:    "nonexistentfile", // Assume this file does not exist
// 			wantErr: true,
// 		},
// 		// Add more test cases here
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := fs.DeleteFile(tt.file)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//
