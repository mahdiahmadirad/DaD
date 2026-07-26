package main

import (
	"os"

	"github.com/mahdiahmadirad/DaD/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, os.Environ()))
}
