package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/Mayank-Creater/go-blog/internal/content"
	"github.com/Mayank-Creater/go-blog/internal/render"
	"github.com/Mayank-Creater/go-blog/internal/server"
	"github.com/fsnotify/fsnotify"
)

//go:embed templates static
var embeddedFiles embed.FS

func helpMenu() {
	fmt.Println("Static Blog Engine CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  ./blog <command> [arguments]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  build       Full rebuild of the static site")
	fmt.Println("  serve       Start the local static HTTP development server with content watching")
	fmt.Println("  new         Create a new post scaffold in content/")
	fmt.Println("\nUse \"./blog <command> -h\" for more information about a command.")
}

func main() {
	if len(os.Args) < 2 {
		helpMenu()
		os.Exit(0)
	}

	homePath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		if err := runBuild(homePath); err != nil {
			fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Build completed successfully.")

	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		portFlag := serveCmd.String("port", ":8080", "The port to run local server")
		if err := serveCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing serve flags: %v\n", err)
			os.Exit(1)
		}

		if err := runBuild(homePath); err != nil {
			fmt.Fprintf(os.Stderr, "Initial build failed: %v\n", err)
			os.Exit(1)
		}

		stopWatch, err := startContentWatcher(homePath, func() {
			fmt.Println("Change detected, rebuilding...")
			if err := runBuild(homePath); err != nil {
				fmt.Fprintf(os.Stderr, "Rebuild failed: %v\n", err)
				return
			}
			fmt.Println("Rebuild completed.")
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error starting watcher: %v\n", err)
			os.Exit(1)
		}
		defer stopWatch()

		server.Serve(*portFlag, embeddedFiles)

	case "new":
		newCmd := flag.NewFlagSet("new", flag.ExitOnError)
		titleFlag := newCmd.String("title", "", "The title of the new post")
		tagFlag := newCmd.String("tag", "new", "The tags of the new post")
		draftFlag := newCmd.Bool("draft", false, "Is the new post draft for now?")
		descriptionFlag := newCmd.String("description", "", "The description of the new post")
		coverFlag := newCmd.String("cover", "", "The cover image of the new post")
		if err := newCmd.Parse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing new flags: %v\n", err)
			os.Exit(1)
		}
		if *titleFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: Missing required flag '-title'")
			fmt.Fprintln(os.Stderr, "Usage: ./blog new -title \"My Post Title\"")
			os.Exit(1)
		}
		if err := content.Scaffold(*titleFlag, *tagFlag, *draftFlag, *descriptionFlag, *coverFlag, homePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating post: %v\n", err)
			os.Exit(1)
		}

	default:
		helpMenu()
		os.Exit(1)
	}
}

func runBuild(homePath string) error {
	contentDir := filepath.Join(homePath, "content")
	if info, err := os.Stat(contentDir); err != nil {
		return fmt.Errorf("content directory error: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("content path is not a directory: %s", contentDir)
	}

	if err := os.MkdirAll(filepath.Join(homePath, "public"), 0755); err != nil {
		return fmt.Errorf("creating public directory: %w", err)
	}

	posts, err := content.ContentScanner(contentDir)
	if err != nil {
		return err
	}
	if err := render.Render(posts, homePath); err != nil {
		return err
	}
	return nil
}

func startContentWatcher(homePath string, onChange func()) (func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	contentDir := filepath.Join(homePath, "content")

	addWatchRecursive := func(root string) error {
		return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if err := watcher.Add(path); err != nil {
					return fmt.Errorf("watch %s: %w", path, err)
				}
			}
			return nil
		})
	}

	if err := addWatchRecursive(contentDir); err != nil {
		watcher.Close()
		return nil, err
	}

	fmt.Printf("Watching for content changes in: %s\n", contentDir)

	var (
		mu      sync.Mutex
		timer   *time.Timer
		buildMu sync.Mutex
	)

	scheduleBuild := func() {
		mu.Lock()
		defer mu.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(200*time.Millisecond, func() {
			buildMu.Lock()
			defer buildMu.Unlock()
			onChange()
		})
	}

	done := make(chan struct{})
	var closeOnce sync.Once

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						_ = addWatchRecursive(event.Name)
					}
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					scheduleBuild()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)

			case <-done:
				return
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		closeOnce.Do(func() {
			close(done)
			_ = watcher.Close()
		})
	}()

	stop := func() {
		closeOnce.Do(func() {
			signal.Stop(sigCh)
			close(done)
			_ = watcher.Close()
		})
		mu.Lock()
		if timer != nil {
			timer.Stop()
		}
		mu.Unlock()
	}

	return stop, nil
}