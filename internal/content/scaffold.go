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

func Scaffold(title string, tags string, draft bool, description string, cover string, homePath string) error {
	filePath := filepath.Join(homePath, "content", fmt.Sprintf("%s.md", slugify(title)))

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("creating content directory: %w", err)
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
	frontmatter = fmt.Sprintf(frontmatter, title, date, tagify(tags), draft, description, cover, title)

	if err := os.WriteFile(filePath, []byte(frontmatter), 0644); err != nil {
		return fmt.Errorf("writing scaffold file %q: %w", filePath, err)
	}

	fmt.Printf("Created new post: %s\n", filePath)
	return nil
}