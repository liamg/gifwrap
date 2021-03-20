package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/liamg/gifwrap/pkg/ascii"

	"github.com/spf13/cobra"
)

var enableFill bool

func main() {

	var rootCmd = &cobra.Command{
		Use:  "gifwrap [url-or-path]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var renderer *ascii.Renderer
			var err error
			arg := args[0]
			if strings.Contains(arg, "://") {
				renderer, err = ascii.FromURL(arg)
			} else {
				renderer, err = ascii.FromFile(arg)
			}

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}

			renderer.SetFill(enableFill)

			if err := renderer.Play(); err != nil && err != ascii.ErrQuit {
				_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
		},
	}
	rootCmd.Flags().BoolVarP(&enableFill, "fill", "f", enableFill, "Fill the entire terminal with the gif, ignoring aspect ratio")
	_ = rootCmd.Execute()
}
