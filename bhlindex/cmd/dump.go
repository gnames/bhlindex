/*
Copyright © 2022 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/io/dbio"
	"github.com/gnames/bhlindex/io/dumpio"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Creates a JSON/CSV/TSV dump of found names",
	Long: `The 'dump' command extracts results of name-detection and
name-verification and writes them to STDOUT. Supports 'json', 'csv',
'tsv' output formats, default format is CSV`,
	Run: func(cmd *cobra.Command, _ []string) {
		f, err := cmd.Flags().GetString("format")
		if f != "" {
			opts = append(opts, config.OptOutputFormat(getFormat(f)))
		}
		cfg := config.New(opts...)
		db := dbio.New(cfg)
		bhli := bhlindex.New(cfg)
		dmp := dumpio.New(cfg, db)
		err = bhli.DumpNames(dmp)
		if err != nil {
			log.Fatal().Err(err).Msg("Name verification failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringP("format", "f", "", "output format: 'csv', 'tsv', 'json'")
}
