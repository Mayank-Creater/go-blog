package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func slugify(text string) string{
	lowerText :=  strings.ToLower(text)
	words := strings.Fields(lowerText)
	return strings.Join(words, "-")
}

func tagify(text string) string{
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return strings.Join(words, ", ")
}

func Scaffold(title string, tags string, draft bool, description string, cover string, homePath string){
	filePath := filepath.Join(homePath, "content", fmt.Sprintf("%s.md", slugify(title)))
	_, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error Creating file: %v\n", err)
	}

	frontmatter := `---
title: "%s"
date: "%s"
tags: [%s]
draft: %t
description: "%s" 
cover: "%s"
---

# %s
`
	now := time.Now()
	date := now.Format("2006-01-02")
	frontmatter = fmt.Sprintf(frontmatter, title, date, tagify(tags), false, description, cover, title)

	err = os.WriteFile(filePath, []byte(frontmatter), 0644)
	if err != nil {
		fmt.Printf("Error Writing File: %v\n", err)
	}
}