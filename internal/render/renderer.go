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

func Render(posts []content.Content){
	fmt.Println("IN render")

	baseTmpl := template.New("base.html").Option("missingkey=error").Funcs(templateFuncs)

	baseTmpl = template.Must(template.ParseFiles("templates/base.html")).Funcs(templateFuncs)
	baseTmpl = template.Must(baseTmpl.ParseGlob("templates/partials/*.html"))

	homeTmpl := template.Must(baseTmpl.Clone())
	homeTmpl = template.Must(homeTmpl.ParseFiles("templates/index.html"))

	var homeBuf bytes.Buffer
	homeFile, _ := os.Create("public/index.html")
	defer homeFile.Close()

	homeData := HomepageData{
		Posts: posts,
	}

	err := homeTmpl.ExecuteTemplate(&homeBuf, "base.html", homeData)
	if err != nil {
		fmt.Println(err)
	}
	_ = os.WriteFile("public/index.html", homeBuf.Bytes(), 0644)


	postTmpl := template.Must(baseTmpl.Clone())
	postTmpl = template.Must(postTmpl.ParseFiles("templates/post.html"))

	targetPost := posts[0]
	
	postFile, _ := os.Create("public/post-hello-world.html")
	defer postFile.Close()
	
	var postBuf bytes.Buffer
	err = postTmpl.ExecuteTemplate(&postBuf, "base.html", targetPost)
	if err != nil {
		fmt.Println(err)
	}
	
	_ = os.WriteFile("public/post.html", postBuf.Bytes(), 0644)

	fmt.Println(postBuf.String())
}