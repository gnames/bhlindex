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

package cmd

import (
	"log"
	"os"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/finder"
	dictionary "github.com/gnames/gnfinder/dict"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Finds name-strings in Biodiversity Heritage Library",
	Long: `Populates an empty database with metadata about pages and volumes
	from Biodiversity Heritage Library. Then it goes volume by volume and
	detects scientific names in them. It uses heuristic and natural language
	processing (Bayes) approaches.`,
	Run: func(cmd *cobra.Command, args []string) {
		workers, err := cmd.Flags().GetInt("workers")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Printf("Processing items with %d workers...", workers)
		db, err := bhlindex.DbInit()
		dict := dictionary.LoadDictionary()
		defer func() {
			e := db.Close()
			bhlindex.Check(e)
		}()
		bhlindex.Check(err)

		finder.ProcessItems(db, dict, workers)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	findCmd.Flags().IntP("workers", "w", 10, "number of name-finding workers")
}
