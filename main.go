package main

import (
	// "fmt"
	"os"
	// "path/filepath"

	"github.com/Mayank-Creater/go-blog/internal/content"
	"github.com/Mayank-Creater/go-blog/internal/render"
)

func main(){
	folderPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	folderPath += "/content"

	result := content.ContentScanner(folderPath)

	render.Render(result)
	
	
}