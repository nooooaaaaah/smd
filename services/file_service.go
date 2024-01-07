package services

import (
	"Smd/types"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type FileService interface {
	ParseAndValidateFile(r *http.Request, w http.ResponseWriter) (multipart.File, *multipart.FileHeader, error)
	HashFile(file multipart.File) (string, error)
	SaveAndUploadFile(file multipart.File, hashedFilename string, subdirectory string, f types.File) error
	UploadFile(f types.File) error
	SaveFile(file multipart.File, hashedFilename, subdirectory string) (string, error)
}

type fileService struct {
	db types.Database
}

func NewFileService() FileService {
	fs := &fileService{
		db: types.NewDatabase(),
	}
	err := fs.db.Connect()
	if err != nil {
		panic(err)
	}

	return fs
}

func (fs *fileService) UploadFile(f types.File) error {
	err := fs.db.InsertFile(f)
	if err != nil {
		return fmt.Errorf("error uploading file to database: %v", err)
	}
	return nil
}

func (fs *fileService) SaveFile(file multipart.File, hashedFilename, subdirectory string) (path string, err error) {
	flatFilename := hashedFilename
	if subdirectory != "" {
		flatFilename = subdirectory + "_" + hashedFilename
	}

	path = filepath.Join("/StoreMeDaddy", flatFilename)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = file.Seek(0, 0); err != nil {
		return "", err
	}

	if _, err = io.Copy(out, file); err != nil {
		return "", err
	}
	if err := out.Close(); err != nil {
		return "", fmt.Errorf("error closing file: %v", err)
	}

	return path, nil
}

func (fs *fileService) SaveAndUploadFile(file multipart.File, hashedFilename string, subdirectory string, f types.File) error {
	filepath, err := fs.SaveFile(file, hashedFilename, subdirectory)
	if err != nil {
		return err
	}
	f.Location = filepath
	err = fs.UploadFile(f)
	if err != nil {
		return err
	}
	return nil
}

func (fs *fileService) ParseAndValidateFile(r *http.Request, w http.ResponseWriter) (file multipart.File, header *multipart.FileHeader, err error) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Max memory 10MB
	if r.Body == nil || r.ContentLength == 0 {
		return nil, nil, fmt.Errorf("request body empty")
	}
	err = r.ParseMultipartForm(10 << 20) // Max memory 10MB

	if err != nil {
		if err.Error() == "http: request body too large" {
			return nil, nil, fmt.Errorf("file too big")
		}
		switch err {
		case http.ErrNotMultipart:
			return nil, nil, fmt.Errorf("request body empty")
		case http.ErrMissingFile:
			return nil, nil, fmt.Errorf("no file in request body")
		case http.ErrHandlerTimeout:
			return nil, nil, fmt.Errorf("request timed out, file might be too big")
		default:
			return nil, nil, fmt.Errorf("error parsing form: %v", err)
		}
	}

	file, header, err = r.FormFile("uploadFile")
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving the file")
	}
	return file, header, nil
}

func (fs *fileService) HashFile(file multipart.File) (hashedFilename string, err error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
