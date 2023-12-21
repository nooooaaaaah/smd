package handlers

import (
	"Smd/services"
	"Smd/types"
	"io"
	"net/http"
	"os"
)

type FileHandler struct {
	fileService services.FileService
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	subdirectory := r.FormValue("subdirectory")

	out, err := os.Create("/StoreMeDaddy/" + subdirectory + header.Filename)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}
	f := types.File{
		Name:        header.Filename,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Location:    "StoreMeDaddy/" + subdirectory + header.Filename,
		OwnerID:     "1",
	}
	fileService := services.NewFileService()
	fileService.UploadFile(f)
	w.Write([]byte("File uploaded successfully"))
}
