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

type HomepageData struct{
	Posts []content.Content
	Title string
}

func Render(posts []content.Content, homePath string){
	fmt.Println("IN render")

	copyStaticFiles(homePath)

	baseTmpl := template.New("base.html").Option("missingkey=error").Funcs(templateFuncs)

	baseTmpl = template.Must(template.ParseFiles("templates/base.html")).Funcs(templateFuncs)
	baseTmpl = template.Must(baseTmpl.ParseGlob("templates/partials/*.html"))

	homeTmpl := template.Must(baseTmpl.Clone())
	homeTmpl = template.Must(homeTmpl.ParseFiles("templates/index.html"))

	var homeBuf bytes.Buffer
	homeFile, _ := os.Create(homePath+"/public/index.html")
	defer homeFile.Close()

	homeData := HomepageData{
		Posts: posts,
		Title: "Go Blog",
	}

	err := homeTmpl.ExecuteTemplate(&homeBuf, "base.html", homeData)
	if err != nil {
		fmt.Println(err)
	}
	_, err = homeBuf.WriteTo(homeFile)
	if err != nil {
		fmt.Println(err)
	}


	postTmpl := template.Must(baseTmpl.Clone())
	postTmpl = template.Must(postTmpl.ParseFiles("templates/post.html"))

	for _, post := range posts {
		filePath := fmt.Sprintf("/public/%s.html", slugify(post.Metadata.Title))
		postFile, _ := os.Create(homePath+filePath)
		defer postFile.Close()
		
		var postBuf bytes.Buffer
		err = postTmpl.ExecuteTemplate(&postBuf, "base.html", post)
		if err != nil {
			fmt.Println(err)
		}
		
		_ = os.WriteFile(homePath+filePath, postBuf.Bytes(), 0644)

	}
}

func copyStaticFiles(homePath string) {
	staticSrc := filepath.Join(homePath, "static")
	staticDst := filepath.Join(homePath, "public", "static")
	
	// Remove existing static folder in public
	os.RemoveAll(staticDst)
	
	err := filepath.Walk(staticSrc, func(path string, info os.FileInfo, err error) error {
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
		
		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		
		_, err = io.Copy(dstFile, srcFile)
		return err
	})
	
	if err != nil {
		fmt.Printf("Error copying static files: %v\n", err)
	}
}