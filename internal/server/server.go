package server

import (
	"fmt"
	"log"
	"net/http"
)

func Serve(port string){
	staticDir := "./public"
	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileServer)

	fmt.Printf("Starting at http://localhost:8080\n")
	fmt.Printf("Serving files from %s\n", staticDir)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}