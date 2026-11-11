package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"io/fs"
)

func Serve(port string, embeddedFiles embed.FS){
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		fmt.Printf("Error Embedding: %v\n", err)
	}
	staticDir := "./public"
	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileServer)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	fmt.Printf("Starting at http://localhost:8080\n")
	fmt.Printf("Serving files from %s\n", staticDir)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}