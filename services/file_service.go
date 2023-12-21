package services

import (
	"Smd/types"

	_ "github.com/mattn/go-sqlite3"
)

type FileService interface {
	UploadFile(f types.File) error
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
		return err
	}
	return nil
}
