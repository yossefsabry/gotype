package main

import (
	"fmt"
	"os"

	"github.com/yossefsabry/gotype/internal/app"
)

func main() {
	// start application
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
