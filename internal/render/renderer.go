package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/Mayank-Creater/go-blog/internal/content"
)

type HomepageData struct {
	Posts []content.Content
	Title string
}

func Render(posts []content.Content, homePath string) error {
	fmt.Println("Rendering site...")

	if err := copyStaticFiles(homePath); err != nil {
		return fmt.Errorf("copying static files: %w", err)
	}

	baseTmpl := template.New("base.html").Option("missingkey=error").Funcs(templateFuncs)
	baseTmpl, err := baseTmpl.ParseFiles(filepath.Join(homePath, "templates", "base.html"))
	if err != nil {
		return fmt.Errorf("parsing base template: %w", err)
	}
	baseTmpl, err = baseTmpl.ParseGlob(filepath.Join(homePath, "templates", "partials", "*.html"))
	if err != nil {
		return fmt.Errorf("parsing partial templates: %w", err)
	}

	homeTmpl, err := baseTmpl.Clone()
	if err != nil {
		return fmt.Errorf("cloning base template for home: %w", err)
	}
	homeTmpl, err = homeTmpl.ParseFiles(filepath.Join(homePath, "templates", "index.html"))
	if err != nil {
		return fmt.Errorf("parsing home template: %w", err)
	}

	homeData := HomepageData{
		Posts: posts,
		Title: "Go Blog",
	}

	var homeBuf bytes.Buffer
	if err := homeTmpl.ExecuteTemplate(&homeBuf, "base.html", homeData); err != nil {
		return fmt.Errorf("executing home template: %w", err)
	}

	homeOutPath := filepath.Join(homePath, "public", "index.html")
	if err := os.MkdirAll(filepath.Dir(homeOutPath), 0755); err != nil {
		return fmt.Errorf("creating public directory: %w", err)
	}
	if err := os.WriteFile(homeOutPath, homeBuf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing home page %q: %w", homeOutPath, err)
	}

	postTmpl, err := baseTmpl.Clone()
	if err != nil {
		return fmt.Errorf("cloning base template for posts: %w", err)
	}
	postTmpl, err = postTmpl.ParseFiles(filepath.Join(homePath, "templates", "post.html"))
	if err != nil {
		return fmt.Errorf("parsing post template: %w", err)
	}

	for _, post := range posts {
		slug := slugify(post.Metadata.Title)
		postOutPath := filepath.Join(homePath, "public", slug+".html")

		var postBuf bytes.Buffer
		if err := postTmpl.ExecuteTemplate(&postBuf, "base.html", post); err != nil {
			return fmt.Errorf("executing post template for %q: %w", post.Metadata.Title, err)
		}

		if err := os.WriteFile(postOutPath, postBuf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing post file %q: %w", postOutPath, err)
		}
	}

	return nil
}

func copyStaticFiles(homePath string) error {
	staticSrc := filepath.Join(homePath, "static")
	staticDst := filepath.Join(homePath, "public", "static")

	if _, err := os.Stat(staticSrc); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("checking static source %q: %w", staticSrc, err)
	}

	if err := os.RemoveAll(staticDst); err != nil {
		return fmt.Errorf("removing existing static destination %q: %w", staticDst, err)
	}

	if err := filepath.Walk(staticSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(staticSrc, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(staticDst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("walking static directory %q: %w", staticSrc, err)
	}

	return nil
}