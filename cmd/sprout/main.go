package main

import (
	"flag"
	"fmt"
	"os"

	"sprout/internal/app"
	"sprout/internal/logx"
)

func main() {
	var (
		root       = flag.String("root", ".", "Site root directory")
		clean      = flag.Bool("clean", false, "Clean public directory before build")
		verbose    = flag.Bool("verbose", false, "Verbose output")
		port       = flag.Int("port", 1313, "Port for serve command")
		rebuild    = flag.String("rebuild", "request", "Rebuild mode: manual, request, watch")
		livereload = flag.Bool("livereload", false, "Enable livereload (dev only)")
	)
	flag.Parse()

	logx.SetVerbose(*verbose)

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: sprout <command> [flags]\n")
		fmt.Fprintf(os.Stderr, "Commands: init, build, serve\n")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "init":
		if err := app.Init(*root); err != nil {
			logx.Errorf("%v", err)
			os.Exit(1)
		}
		fmt.Printf("Initialized Sprout site in %s\n", *root)
	case "build":
		result, err := app.Build(*root, *clean)
		if err != nil {
			logx.Errorf("%v", err)
			os.Exit(1)
		}
		fmt.Printf("Build complete: %d pages, %d assets\n", result.BuiltPages, result.CopiedAssets)
	case "serve":
		rebuildMode := app.RebuildMode(*rebuild)
		if err := app.Serve(*root, *port, rebuildMode, *livereload); err != nil {
			logx.Errorf("%v", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintf(os.Stderr, "Commands: init, build, serve\n")
		os.Exit(1)
	}
}
