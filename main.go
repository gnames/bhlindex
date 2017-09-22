package main

import (
	"fmt"
	"os"
)

var githash string = "n/a"
var buildstamp string = "n/a"

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	switch command {
	case "version":
		fmt.Printf(" Git commit hash: %s\n UTC Build Time: %s\n\n",
			githash, buildstamp)
	case "index":
		fmt.Println("Indexing...")
	default:
		help := `
Usage:
  bhlindex version
  bhlindex index
`
		fmt.Println(help)
	}
}
