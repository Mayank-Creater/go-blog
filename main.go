package main

import (
	"flag"
	"fmt"
	"os"
	"embed"

	"github.com/Mayank-Creater/go-blog/internal/content"
	"github.com/Mayank-Creater/go-blog/internal/render"
	"github.com/Mayank-Creater/go-blog/internal/server"
)

//go:embed templates static
var embeddedFiles embed.FS

func helpMenu(){
	fmt.Println("Static Blog Engine CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  ./blog <command> [arguments]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  build       Full rebuild of the static site")
	fmt.Println("  serve       Start the local static HTTP development server")
	fmt.Println("\nUse \"./blog <command> -help\" for more information about a command.")
}

func main(){
	if len(os.Args) < 2 {
		helpMenu()
		os.Exit(0)
	}

	homePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	titleFlag := newCmd.String("title", "", "The title of the new post")
	tagFlag := newCmd.String("tag", "new", "The tags of the new post")
	draftFlag := newCmd.Bool("draft", false, "Is the new post draft for now?")
	descriptionFlag := newCmd.String("description", "", "The description of the new post")
	coverFlag := newCmd.String("cover", "", "The cover image of the new post")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag := serveCmd.String("port", ":8080", "The port to run local server")

	switch os.Args[1]{
	case "serve":
		server.Serve(*portFlag, embeddedFiles)
	case "build":
		folderPath := homePath + "/content"
		result := content.ContentScanner(folderPath)
		render.Render(result, homePath)
	case "new":
		newCmd.Parse(os.Args[2:])
		if *titleFlag == "" {
			fmt.Println("Error: Missing required flag '-title'")
			fmt.Println("Usage: ./blog generate -title \"My Post Title\"")
			os.Exit(1)
		}
		content.Scaffold(*titleFlag, *tagFlag, *draftFlag, *descriptionFlag, *coverFlag, homePath)
	}
	
}