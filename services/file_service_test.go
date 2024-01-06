package services

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseAndValidateFile(t *testing.T) {
	fs := &fileService{}

	tests := []struct {
		name    string
		body    []byte
		wantErr bool
		errMsg  string
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
			wantErr: true,
			errMsg:  "request body empty",
		},
		{
			name:    "NilBody",
			body:    nil,
			wantErr: true,
			errMsg:  "request body empty",
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest("POST", "/upload", bytes.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
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
