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
	"github.com/gnames/bhlindex/internal"
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/io/dbio"
	"github.com/gnames/bhlindex/internal/io/dumpio"
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
		var err error

		f, _ := cmd.Flags().GetString("format")
		if f != "" {
			opts = append(opts, config.OptOutputFormat(getFormat(f)))
		}
		d, _ := cmd.Flags().GetString("dir")
		if d != "" {
			d, err = getDir(d)
			if err != nil {
				log.Fatal().Err(err).Msgf("Cannot use directory '%s'", d)
			}
			opts = append(opts, config.OptOutputDir(d))
		}
		s, _ := cmd.Flags().GetString("sources")
		if s != "" {
			var ss []int
			ss, err = getSources(s)
			if err != nil {
				log.Fatal().Err(err).Msgf("Cannot parse data-sources '%s'", s)
			}
			opts = append(opts, config.OptOutputDataSourceIDs(ss))
		}

		cfg := config.New(opts...)
		db := dbio.New(cfg)
		bhli := bhlindex.New(cfg)
		dmp := dumpio.New(cfg, db)

		err = bhli.DumpPages(dmp)
		if err != nil {
			log.Fatal().Err(err).Msg("Dump of pages failed")
		}

		err = bhli.DumpNames(dmp)
		if err != nil {
			log.Fatal().Err(err).Msg("Dump of names failed")
		}

		err = bhli.DumpOccurrences(dmp)
		if err != nil {
			log.Fatal().Err(err).Msg("Dump of occurrences failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringP("format", "f", "",
		"output format: 'csv', 'tsv', 'json', default 'csv'")
	dumpCmd.Flags().StringP("dir", "d", "",
		"output directory, defult current dir")
	dumpCmd.Flags().StringP("sources", "s", "",
		"filter by gnverifier data-sources. Example: 1,11")
}
