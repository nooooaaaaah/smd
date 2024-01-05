package handlers

import (
	"Smd/services"
	"Smd/types"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileHandler struct {
	FileService services.FileService
}

func (fh *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := parseAndValidateFile(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	hashedFilename, err := hashFile(file)
	if err != nil {
		http.Error(w, "Error hashing the file", http.StatusInternalServerError)
		return
	}

	filepath, err := saveFile(file, hashedFilename, r.FormValue("subdirectory"))
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	f := types.File{
		Name:        header.Filename,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Location:    filepath,
		OwnerID:     "1",
	}
	if err := fh.FileService.UploadFile(f); err != nil {
		http.Error(w, "Error MetaData not saved", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func parseAndValidateFile(r *http.Request, w http.ResponseWriter) (file multipart.File, header *multipart.FileHeader, err error) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Max memory 10MB
	err = r.ParseMultipartForm(10 << 20)            // Max memory 10MB
	if err != nil {
		return nil, nil, fmt.Errorf("file too big")
	}

	file, header, err = r.FormFile("uploadFile")
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving the file")
	}

	return file, header, nil
}

func hashFile(file multipart.File) (hashedFilename string, err error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hashedFilename = hex.EncodeToString(hasher.Sum(nil))

	return hashedFilename, nil
}

func saveFile(file multipart.File, hashedFilename, subdirectory string) (path string, err error) {
	flatFilename := subdirectory + "_" + hashedFilename

	path = filepath.Join("/StoreMeDaddy", flatFilename)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = file.Seek(0, 0) // Reset file read position
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(out, file); err != nil {
		return "", err
	}

	return path, nil
}
