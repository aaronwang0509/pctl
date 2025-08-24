package main

import (
	"os"

	"github.com/aaronwang/pctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}