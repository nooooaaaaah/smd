package services

import (
	"Smd/types"
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
		name    string
		file    multipart.File
		subDir  string
		wantErr bool
	}{
		{
			name:   "Success",
			subDir: "",
			file: func() multipart.File {
				file, err := os.CreateTemp("", "test")
				if err != nil {
					t.Fatal(err)
				}
				file.WriteString("test data")
				file.Seek(0, 0)
				return file
			}(), // Simulate a file with "test data"
			wantErr: false,
		},
		{
			name:    "NilFile",
			subDir:  "",
			file:    nil,
			wantErr: true,
		},
		{
			name:    "Subdirectory",
			subDir:  "taco",
			file:    func() multipart.File { file, _ := os.CreateTemp("", "test"); return file }(), // Simulate an empty file
			wantErr: false,
			// Add more test cases here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.SaveFile(tt.file, tt.name, tt.subDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewFileService(t *testing.T) {
	fs := NewFileService()

	if fs == nil {
		t.Errorf("NewFileService() = %v, want non-nil", fs)
	}
}

func TestUploadFile(t *testing.T) {
	fs := NewFileService()

	tests := []struct {
		name    string
		file    types.File
		wantErr bool
	}{
		{
			name:    "Success",
			file:    types.File{ID: "1", Name: "testfile", Location: "/tmp/testfile"},
			wantErr: false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Remove("smd.db")
			types.Database.CreateDb(types.NewDatabase())
			err := fs.UploadFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveAndUploadFile(t *testing.T) {
	fs := NewFileService()

	tests := []struct {
		name           string
		file           multipart.File
		hashedFilename string
		subdirectory   string
		f              types.File
		wantErr        bool
	}{
		{
			name:           "Success",
			file:           func() multipart.File { file, _ := os.CreateTemp("", "test"); return file }(), // Simulate an empty file
			hashedFilename: "hashedFilename",
			subdirectory:   "",
			f:              types.File{ID: "2", Name: "SaveAndUploadFile", Location: "/tmp/testfile"},
			wantErr:        false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			types.Database.CreateDb(types.NewDatabase())
			err := fs.SaveAndUploadFile(tt.file, tt.hashedFilename, tt.subdirectory, tt.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveAndUploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashingFile(t *testing.T) {
	fs := NewFileService()

	tests := []struct {
		file    multipart.File
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			file: func() multipart.File {
				file, _ := os.CreateTemp("", "test")
				file.WriteString("test data")
				file.Seek(0, 0)
				return file
			}(), // Simulate a file with "test data"
			wantErr: false,
		},
		{
			name:    "NilFile",
			file:    nil,
			wantErr: true,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.HashFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	fs := NewFileService()
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "Success",
			fileName: "1",
			wantErr:  false,
		},
		{
			name:     "NonExistentFile",
			fileName: "2",
			wantErr:  true,
		},
		// Add more test cases here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.GetFile(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
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
