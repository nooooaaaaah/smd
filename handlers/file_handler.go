package handlers

import (
	"Smd/services"
	"Smd/types"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type FileHandler struct {
	fileService services.FileService
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Max memory 10MB
	err := r.ParseMultipartForm(10 << 20)           // Max memory 10MB
	if err != nil {
		http.Error(w, "File to big", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		http.Error(w, "Error hashing the file", http.StatusInternalServerError)
		return
	}
	hashedFilename := hex.EncodeToString(hasher.Sum(nil))
	subdirectory := r.FormValue("subdirectory")
	flatFilename := subdirectory + "_" + hashedFilename
	filepath := filepath.Join("/StoreMeDaddy", flatFilename)
	out, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = file.Seek(0, 0) // Reset file read position
	if err != nil {
		http.Error(w, "Error resetting file position", http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(out, file); err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}
	f := types.File{
		Name:        header.Filename,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Location:    filepath,
		OwnerID:     "1",
	}
	fileService := services.NewFileService()
	if err := fileService.UploadFile(f); err != nil {
		http.Error(w, "Error MetaData not saved", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
