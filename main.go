package main

import (
	handlers "Smd/handlers"
	"Smd/services"
	"Smd/types"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	maxFileSize := os.Getenv("MAX_FILE_SIZE")
	if port == "" {
		port = "8080"
	}
	if maxFileSize == "" {
		maxFileSize = "100000000"
	}

	fileService := services.NewFileService()
	fileHandler := &handlers.FileHandler{
		FileService: fileService,
	}
	fmt.Println("Creating ~/StoreMeDaddy directory")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(homeDir, "StoreMeDaddy")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting server (modem noises)...")
	fmt.Println("Registering handler for /upload")
	http.HandleFunc("/upload", fileHandler.UploadFileHandler)
	fmt.Println("Handlers registered")
	fmt.Println("Spinning up database")
	db := types.NewDatabase()
	types.Database.CreateDb(db)
	fmt.Println("Server started")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println("Listening on port " + port)
	if err := http.ListenAndServe(":5464", nil); err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}
