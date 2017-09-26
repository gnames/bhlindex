package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GlobalNamesArchitecture/bhlindex"
	"github.com/GlobalNamesArchitecture/bhlindex/finder"
)

var githash = "n/a"
var buildstamp = "n/a"

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
		makeIndex()
	default:
		help := `
Usage:
  bhlindex version
  bhlindex index
`
		fmt.Println(help)
	}
}

func makeIndex() {
	log.Println("Processing titles...")
	db, err := bhlindex.DbInit()
	defer func() {
		e := db.Close()
		bhlindex.Check(e)
	}()
	bhlindex.Check(err)

	finder.ProcessTitles(db)
}
