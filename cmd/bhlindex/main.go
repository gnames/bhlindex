// Copyright © 2018 Dmitry Mozzherin <dmozzherin@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/finder"
	dictionary "github.com/gnames/gnfinder/dict"
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
	dict := dictionary.LoadDictionary()
	defer func() {
		e := db.Close()
		bhlindex.Check(e)
	}()
	bhlindex.Check(err)

	finder.ProcessTitles(db, &dict)
}
