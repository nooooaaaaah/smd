package handlers

import (
	"Smd/types"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) ParseAndValidateFile(r *http.Request, w http.ResponseWriter) (multipart.File, *multipart.FileHeader, error) {
	args := m.Called(r, w)
	return args.Get(0).(multipart.File), args.Get(1).(*multipart.FileHeader), args.Error(2)
}

func (m *MockFileService) HashFile(file multipart.File) (string, error) {
	args := m.Called(file)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) SaveAndUploadFile(file multipart.File, hashedFilename string, subdirectory string, f types.File) error {
	args := m.Called(file, hashedFilename, subdirectory, f)
	return args.Error(0)
}

func (m *MockFileService) UploadFile(f types.File) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *MockFileService) SaveFile(file multipart.File, hashedFilename, subdirectory string) (string, error) {
	args := m.Called(file, hashedFilename, subdirectory)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) GetFile(id string) (types.File, error) {
	args := m.Called(id)
	return args.Get(0).(types.File), args.Error(1)
}

type FakeMultiPart struct {
	*bytes.Buffer
}

func (fc *FakeMultiPart) Close() error {
	return nil
}

func (fc *FakeMultiPart) ReadAt(p []byte, off int64) (n int, err error) {
	// Convert the offset to an int
	offset := int(off)

	// Check if the offset is out of bounds
	if offset < 0 || offset >= fc.Len() {
		return 0, io.EOF
	}

	// Copy the data from the buffer to p
	n = copy(p, fc.Bytes()[offset:])

	return n, nil
}

func (fc *FakeMultiPart) Seek(offset int64, whence int) (int64, error) {
	// Convert the offset to an int
	o := int(offset)

	// Check if the offset is out of bounds
	if o < 0 || o >= fc.Len() {
		return 0, io.EOF
	}

	// Seek the buffer
	fc.Buffer = bytes.NewBuffer(fc.Bytes()[o:])

	return offset, nil
}

func TestUploadFileHandler(t *testing.T) {
	testCases := []struct {
		name           string
		fileData       []byte
		expectedStatus int
	}{
		{
			name:           "FileTooLarge",
			fileData:       make([]byte, 10*1024*1024+10), // 10 MB + 10 byte
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "Success",
			fileData:       []byte("test data"), // Less than 10MB
			expectedStatus: http.StatusOK,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockFileService)
			fh := FileHandler{
				FileService: mockService,
				StorePath:   "/tmp",
				MaxFileSize: 10 * 1024 * 1024, // 10 MB
			}
			if int64(len(tc.fileData)) > fh.MaxFileSize {
				mockService.On("ParseAndValidateFile", mock.Anything, mock.Anything).Return(&FakeMultiPart{bytes.NewBuffer(tc.fileData)}, &multipart.FileHeader{}, errors.New("file too large"))
			} else {
				mockService.On("ParseAndValidateFile", mock.Anything, mock.Anything).Return(&FakeMultiPart{bytes.NewBuffer(tc.fileData)}, &multipart.FileHeader{}, nil)
			}
			mockService.On("HashFile", mock.Anything).Return("test", nil)
			mockService.On("SaveAndUploadFile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockService.On("UploadFile", mock.Anything).Return(nil)
			mockService.On("SaveFile", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
			// Create a buffer to hold the form data
			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			// Create a file field in the form
			fw, err := w.CreateFormFile("file", "test.txt")
			if err != nil {
				t.Fatal(err)
			}

			// Write data to the file field
			_, err = fw.Write(tc.fileData)
			if err != nil {
				t.Fatal(err)
			}

			// Close the multipart writer or the form will be missing the ending boundary
			err = w.Close()
			if err != nil {
				t.Fatal(err)
			}

			// Create a new request with the form data
			req, err := http.NewRequest("POST", "/upload", &b)
			if err != nil {
				t.Fatal(err)
			}

			// Set the content type to multipart/form-data and include the boundary
			req.Header.Set("Content-Type", w.FormDataContentType())

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the UploadFileHandler
			fh.UploadFileHandler(rr, req)

			// Check the status code
			fmt.Println(rr.Code)
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}
		})
	}
}
