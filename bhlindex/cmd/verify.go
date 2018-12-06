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
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Checks found name-strings against gnindex database",
	Long: `Verifies scientific name-strings found in BHL against
	https://index.globalnames.org service. Results indicate if
	the name-strings are found in a variety of biodiversity databases as
	exact or fuzzy matches.`,
	Run: func(cmd *cobra.Command, args []string) {
		workersNum, err := cmd.Flags().GetInt("workers")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Printf("Verifying name-strings with %d workers...", workersNum)
		db, err := bhlindex.DbInit()
		defer func() {
			e := db.Close()
			bhlindex.Check(e)
		}()
		bhlindex.Check(err)

		finder.Verify(db, workersNum)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	verifyCmd.Flags().IntP("workers", "w", 10, "number of name-finding workers")
}
