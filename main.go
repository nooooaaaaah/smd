package smd

import (
	handlers "Smd/handlers"
	"Smd/types"
	"fmt"
	"net/http"
	"os"
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

	fmt.Println("Starting server...")
	fmt.Println("Registering handler for /upload")
	http.HandleFunc("/upload", handlers.UploadFileHandler)
	fmt.Println("Handlers registered")
	http.ListenAndServe(":8080", nil)
	db := types.NewDatabase()
	types.Database.CreateDb(db)
	fmt.Println("Server started at port 8080")
	fmt.Println("Press Ctrl+C to exit")
}
