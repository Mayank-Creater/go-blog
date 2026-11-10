package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Mayank-Creater/go-blog/internal/content"
	"github.com/Mayank-Creater/go-blog/internal/render"
	"github.com/Mayank-Creater/go-blog/internal/server"
)

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

	// buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag := serveCmd.String("port", ":8080", "The port to run local server")

	switch os.Args[1]{
	case "serve":
		port := *portFlag
		server.Serve(port)
	case "build":
		homePath, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		folderPath := homePath + "/content"

		result := content.ContentScanner(folderPath)

		render.Render(result, homePath)
		
	}
	
}