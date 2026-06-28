package main

import (
	"flag"

	"github.com/anton-fuji/gitviz/internal"
)

func main() {
	var folder string
	var graphEmail string
	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
	flag.StringVar(&graphEmail, "graph", "your@email.com", "the email to scan")
	flag.Parse()

	if folder != "" {
		gitviz.Scan(folder)
		return
	}
	gitviz.Stats(graphEmail)
}
