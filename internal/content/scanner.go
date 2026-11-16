package content

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func ContentScanner(folderPath string) ([]Content, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("reading content directory %q: %w", folderPath, err)
	}

	var contentList []Content

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		fullPath := filepath.Join(folderPath, entry.Name())
		fmt.Printf("Reading file: %s\n", entry.Name())

		input, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", fullPath, err)
			continue
		}

		metadata, postContent, err := ContentParser(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", fullPath, err)
			continue
		}

		if metadata.Draft {
			continue
		}

		contentList = append(contentList, Content{
			Metadata: metadata,
			Content:  postContent,
		})
	}

	slices.SortFunc(contentList, func(a, b Content) int {
		return cmp.Compare(a.Metadata.Date, b.Metadata.Date)
	})

	return contentList, nil
}