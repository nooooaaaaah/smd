package handlers

import (
	"Smd/services"
	"Smd/types"
	"net/http"
)

type FileHandler struct {
	FileService services.FileService
	StorePath   string // Configurable store path
	MaxFileSize int64  // Configurable max file MaxFileSize
}

func (fh *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := fh.FileService.ParseAndValidateFile(r, w)
	if err != nil {
		if err.Error() == "file too large" {
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	hashedFilename, err := fh.FileService.HashFile(file)
	if err != nil {
		http.Error(w, "Error hashing the file", http.StatusInternalServerError)
		return
	}

	f := types.File{
		Name:        header.Filename,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Location:    hashedFilename,
		OwnerID:     "1", // todo implement auth
	}

	if err := fh.FileService.SaveAndUploadFile(file, hashedFilename, r.FormValue("subdirectory"), f); err != nil {
		http.Error(w, "Error saving and uploading the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
