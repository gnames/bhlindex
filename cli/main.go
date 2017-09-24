package main

import (
	"fmt"
	"os"

	"github.com/GlobalNamesArchitecture/bhlindex/finder"
	"github.com/GlobalNamesArchitecture/bhlindex/util"
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
	db, err := util.DbInit()
	defer func() {
		e := db.Close()
		util.Check(e)
	}()
	util.Check(err)

	fmt.Println("Collecting titles to the database...")
	finder.FindNames(db)
}
