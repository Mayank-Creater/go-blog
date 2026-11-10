package render

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/Mayank-Creater/go-blog/internal/content"
)

type HomepageData struct{
	Posts []content.Content
}

func Render(posts []content.Content, homePath string){
	fmt.Println("IN render")

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