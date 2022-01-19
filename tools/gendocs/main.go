package main

import (
	"os"

	"github.com/hashicorp/nomad-pack/internal/cli"
)

func main() {
	cli.ExposeDocs = true
	os.Exit(cli.Main([]string{"nomad-pack", "cli-docs"}))
}
