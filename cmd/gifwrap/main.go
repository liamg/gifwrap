package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/liamg/gifwrap/pkg/ascii"
)

func main() {

	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, `Usage: 
	gifwrap [url]
	gifwrap [path]
`)
		os.Exit(1)
	}

	var renderer *ascii.Renderer
	var err error
	arg := os.Args[1]
	if strings.Contains(arg, "://") {
		renderer, err = ascii.FromURL(arg)
	} else {
		renderer, err = ascii.FromFile(arg)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	if err := renderer.Play(); err != nil && err != ascii.ErrQuit {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
