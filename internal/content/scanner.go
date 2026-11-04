package content

import (
	"fmt"
	"os"
	"path/filepath"
)

func ContentScanner(folderPath string) []Content{
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		panic(err)
	}

	var content []Content

	for _, entry := range entries{
		if entry.IsDir(){
			continue
		}
		if filepath.Ext(entry.Name()) != ".md"{
			continue
		}

		fullPath := filepath.Join(folderPath, entry.Name())
		fmt.Printf("Reading file: %s\n", entry.Name())

		input, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Println("Error Reading file:", err)
			continue
		}

		metadata, postContent, err := ContentParser(input)
		if err != nil {
			fmt.Println("Error parsing file:", err)
			continue
		}
	
		content = append(content, Content{
			Metadata: metadata,
			Content: postContent,
		})
	}	
	return content
}