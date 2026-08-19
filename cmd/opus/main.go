// Command opus is the entry point for the opus CLI.
package main

import (
	"os"

	"opus/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
